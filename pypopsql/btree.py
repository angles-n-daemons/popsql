from enum import Enum

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
    ):
        self.data = data
        self.page_size = len(data)
        
        node_type_bytes = data[0:1]
        self.node_type = NodeType(b2i(node_type_bytes))

        num_cells_bytes = data[3:5]
        self.num_cells = b2i(num_cells_bytes)

        cell_offset_bytes = data[5:7]
        self.cell_offset = b2i(cell_offset_bytes)

        self.cells = self.read_cells()

        self.right_pointer = None
        if not self.is_leaf():
            right_pointer_bytes = data[8:12]
            self.right_pointer = b2i(right_pointer_bytes)

        """
        The following fields (first_freeblock, num_fragmented_bytes) are omitted from
        usage, but included for testing to assert that the parser is working correctly
        """

        first_freeblock_bytes = data[1:3]
        self.first_freeblock = b2i(first_freeblock_bytes)

        self.num_fragmented_bytes = data[7]

    def read_cells(self):
        hdrlen = 8 if self.is_leaf() else 12

        cells = []

        for i in range(self.num_cells):
            offset = hdrlen + (i * 2)
            p = b2i(self.data[offset:offset + 2])
            cell = TableLeafCell(self.data, p)
            cells.append(cell)

        return cells

    def is_leaf(self):
        return self.node_type in (NodeType.TABLE_LEAF, NodeType.INDEX_LEAF)

    def _debug_print_cells(self):
        for cell in self.cells:
            print('row id: ', cell.row_id)
            print('payload size: ', cell.payload_size)
            print('payload: ', cell.payload.hex())
            print('cursor end: ', cell.cursor)
            print('\n\n')

class TableLeafCell:
    def __init__(
        self,
        data: bytes,
        cursor: int,
    ):
        self.payload_size, num_read = varint(data, cursor)
        cursor += num_read

        self.row_id, num_read = varint(data, cursor)
        cursor += num_read

        # TODO: does not address overflow
        self.payload = data[cursor:cursor+self.payload_size]
        self.cursor = cursor + self.payload_size
