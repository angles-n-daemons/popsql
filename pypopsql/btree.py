from enum import Enum
from typing import Tuple, List

from record import Record
from util import b2i, varint

class NodeType(Enum):
    INDEX_INTERIOR = 2
    TABLE_INTERIOR = 5
    INDEX_LEAF = 10
    TABLE_LEAF = 13

class Node:
    def __init__(
        self,
        data: bytes,
        db_header: bool = False,
    ):
        offset = 100 if db_header else 0

        self.data = data
        self.db_header = db_header
        self.page_size = len(data)
        
        self.node_type, \
        self.cell_offset, \
        self.num_cells, \
        self.right_pointer, \
        self.first_freeblock, \
        self.num_fragmented_bytes = self.read_header_bytes(data, db_header)

        self.cells = self.read_cells()


    def read_header_bytes(
        self,
        data: bytes,
        db_header: bool,
    ) -> Tuple[
        NodeType,
        int, # num cells
        int, # cell_offset
        int, # right_pointer
        int, # first_freeblock
        int, # num_fragmented_bytes
    ]:
        offset = 100 if db_header else 0

        node_type = NodeType(b2i(data[offset + 0: offset + 1]))

        num_cells = b2i(data[offset + 3: offset + 5])

        cell_offset = b2i(data[offset + 5: offset + 7])

        right_pointer = None
        if not self.is_leaf(node_type):
            right_pointer_bytes = data[offset + 8: offset + 12]
            right_pointer = b2i(right_pointer_bytes)

        """
        The following fields (first_freeblock, num_fragmented_bytes) are omitted from
        usage, but included for testing to assert that the parser is working correctly
        """

        first_freeblock = b2i(data[offset + 1: offset + 3])
        num_fragmented_bytes = data[offset + 7]

        return (
            node_type,
            cell_offset,
            num_cells,
            right_pointer,
            first_freeblock,
            num_fragmented_bytes,
        )

    def read_cells(self) -> List[any]:
        page_header_len = 8 if self.is_leaf() else 12
        db_header_len = 100 if self.db_header else 0

        cells = []

        for i in range(self.num_cells):
            offset = page_header_len + db_header_len + (i * 2)
            p = b2i(self.data[offset:offset + 2])
            cell = TableLeafCell(self.data, p)
            cells.append(cell)

        return cells

    def is_leaf(self, node_type: NodeType = None):
        node_type = node_type or self.node_type
        return node_type in (NodeType.TABLE_LEAF, NodeType.INDEX_LEAF)

    def _debug_print_cells(self):
        for cell in self.cells:
            print('row id: ', cell.row_id)
            print('payload size: ', cell.payload_size)
            print('payload: ', cell.payload.hex())
            print('cursor end: ', cell.cursor)
            cell.record._debug_print_values()
            print('\n\n')

class TableLeafCell:
    def __init__(
        self,
        data: bytes,
        cursor: int,
    ):
        self.payload_size, cursor = varint(data, cursor)

        self.row_id, cursor = varint(data, cursor)

        # TODO: does not address overflow
        self.payload = data[cursor:cursor+self.payload_size]
        self.record = Record(data, cursor)

        self.cursor = cursor + self.payload_size
