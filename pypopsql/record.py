from util import varint

class Record:
    def __init__(
        self,
        data: bytes,
        cursor: int,
    ):
        self.data = data
        self.cursor = cursor

        self.column_types, cursor = self.read_column_types(data, cursor)
        self.values, cursor = self.read_values(data, cursor)

    def read_column_types(
        self,
        data: bytes,
        cursor: int,
    ):
        # The header begins with a single varint which determines the total number of bytes in the header
        # the varint value is the size of the header including the size varint itself
        cursor_start = cursor
        num_bytes_header, num_bytes = varint(data, cursor)
        print('header size bytes', num_bytes_header)
        return [], 0

    def read_values(
        self,
        data: bytes,
        cursor: int,
    ):
        return [], 0

    def _debug_print_values(self):
        print('record')
