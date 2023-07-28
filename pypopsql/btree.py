from enum import Enum
from typing import Tuple

def b2i(b: bytes) -> int:
    return int.from_bytes(b, 'big', signed=False)

def varint(b: bytes, start: int) -> Tuple[int, int]:
    """
    varint reads a variable length int from some byte string
    it takes the following parameters
     - b: sequence of bytes, generally a full table page
     - i: starting index of the varint

    and returns a tuple including the value of the integer, and the
    number of bytes read.
    """

    result = 0

    for j in range(8):
        i = start + j

        # read the next byte into an integer
        byte_num = b[i]

        # result shifts left 7 bits, then the first 7 bits of byte_num are appended
        result = (result << 7) | (byte_num & 0x7f)

        # check the first bit of byte_num to see if we should continue reading
        continue_reading = byte_num & 0x80

        if not continue_reading:
            return result, j+1

    # read last byte, use all 8 bytes to fill the remaining spaces
    byte_num = b[start + 8]
    result = (result << 8) | byte_num

    return result, 9

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

    def test_read_cells(self):
        #TODO calculate actual starting offset
        hdrlen =  8

        cells = []
        for i in range(self.num_cells):
            offset = hdrlen + (i * 2)
            p = b2i(self.data[offset:offset + 2])
            print(f'cell {i}, offset {p}')

class TableLeafCell:
    def __init__(
        self,
        data: bytes,
        start: int,
    ):
        self.payload_size = None
        self.row_id = None
        self.payload = None
