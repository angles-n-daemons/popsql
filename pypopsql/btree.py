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

    def read_cells(self):
        #TODO calculate actual starting offset
        hdrlen =  8

        cells = []

        for i in range(self.num_cells):
            offset = hdrlen + (i * 2)
            p = b2i(self.data[offset:offset + 2])
            cell = TableLeafCell(self.data, p)
            cells.append(cell)

        return cells

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
