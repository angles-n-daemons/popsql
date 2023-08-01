from enum import Enum

from util import varint

class ColumnType(Enum):
    NULL = 0
    TINYINT = 1
    SMALLINT = 2
    SMALLISHINT = 3
    INTEGER = 4
    BIGGISHINT = 5
    LONG = 6
    IEEE754INT = 7
    ZERO = 8
    ONE = 9
    RESERVED_1 = 10
    RESERVED_2 = 11
    BLOB = 12
    TEXT = 13

    def __init__(self, value):
        self._value_ = value
        self.length = None

    @classmethod
    def from_varint(cls, value: int):
        if value < 12:
            return cls(value)

        is_even = value % 2 == 0
        if is_even:
            column_type = cls(12)
            column_type.length = (value-12) // 2
            return column_type
        else:
            column_type = cls(13)
            column_type.length = (value-13) // 2
            return column_type

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
        column_types = []
        cursor_start = cursor
        num_bytes_header, cursor = varint(data, cursor)

        while cursor - cursor_start < num_bytes_header:
            column_type, cursor = varint(data, cursor)

        print('header size bytes', num_bytes_header)
        return column_types, 0

    def read_values(
        self,
        data: bytes,
        cursor: int,
    ):
        return [], 0

    def _debug_print_values(self):
        print('record')
