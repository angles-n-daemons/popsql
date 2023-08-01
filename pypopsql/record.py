from enum import Enum
from typing import Tuple

from util import b2i, varint

class ColumnType(Enum):
    UNKNOWN = -1
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
        column_types = []
        cursor_start = cursor
        num_bytes_header, cursor = varint(data, cursor)

        while cursor - cursor_start < num_bytes_header:
            column_type_int, cursor = varint(data, cursor)
            column_types.append(ColumnType.from_varint(column_type_int))

        return column_types, cursor

    def read_values(
        self,
        data: bytes,
        cursor: int,
    ):
        values = []
        for column_type in self.column_types:
            value, cursor = slelf.read_value(column_type, data, cursor)
            values.append(value)
        return values, cursor
    
    @staticmethod
    def read_value(
        column_type: ColumnType,
        data: bytes,
        cursor: int,
    ) -> Tuple[any, int]:
        if column_type == ColumnType.NULL:
            return None, cursor
        elif column_type == ColumnType.TINYINT:
            return int(data[cursor]), cursor + 1
        elif column_type == ColumnType.SMALLINT:
            return b2i(data[cursor: cursor + 2]), cursor + 2
        elif column_type == ColumnType.SMALLISHINT:
            return b2i(data[cursor: cursor + 3]), cursor + 3
        elif column_type == ColumnType.INTEGER:
            return b2i(data[cursor: cursor + 4]), cursor + 4
        elif column_type == ColumnType.BIGGISHINT:
            return b2i(data[cursor: cursor + 6]), cursor + 6
        elif column_type == ColumnType.LONG:
            return b2i(data[cursor: cursor + 8]), cursor + 8
        elif column_type == ColumnType.ZERO:
            return 0, cursor
        elif column_type == ColumnType.ONE:
            return 1, cursor
        else:
            raise Exception(f'cannot parse column type {column_type}')

    def _debug_print_values(self):
        for column in self.column_types:
            print('column type', column)
