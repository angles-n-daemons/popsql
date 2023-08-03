from dataclasses import dataclass
from typing import Union
from unittest import TestCase

from record import Record, ColumnType, Column

@dataclass
class ColumnTestCase:
    value: int

    expected: ColumnType
    expected_length: Union[int, None]
    throws_error: bool

@dataclass
class ReadValueTestCase:
    column: Column
    data: bytes
    cursor_start: int

    expected: any
    expected_cursor_end: int
    throws_error: bool

class TestRecord(TestCase):
    def test_column_int_parsing(self):
        tests = [
            # test 1-11
            ColumnTestCase(0, ColumnType.NULL, None, False),
            ColumnTestCase(1, ColumnType.TINYINT, None, False),
            ColumnTestCase(2, ColumnType.SMALLINT, None, False),
            ColumnTestCase(3, ColumnType.SMALLISHINT, None, False),
            ColumnTestCase(4, ColumnType.INTEGER, None, False),
            ColumnTestCase(5, ColumnType.BIGGISHINT, None, False),
            ColumnTestCase(6, ColumnType.LONG, None, False),
            ColumnTestCase(7, ColumnType.IEEE754INT, None, False),
            ColumnTestCase(8, ColumnType.ZERO, None, False),
            ColumnTestCase(9, ColumnType.ONE, None, False),
            ColumnTestCase(10, ColumnType.RESERVED_1, None, False),
            ColumnTestCase(11, ColumnType.RESERVED_2, None, False),

            # test less than 1
            ColumnTestCase(-1, None, None, True),
            ColumnTestCase(-10, None, None, True),

            # test blob w multiple values
            ColumnTestCase(12, ColumnType.BLOB, 0, False),
            ColumnTestCase(14, ColumnType.BLOB, 1, False),
            ColumnTestCase(16, ColumnType.BLOB, 2, False),
            ColumnTestCase(140, ColumnType.BLOB, 64, False),

            # test text with multiple values
            ColumnTestCase(13, ColumnType.TEXT, 0, False),
            ColumnTestCase(15, ColumnType.TEXT, 1, False),
            ColumnTestCase(17, ColumnType.TEXT, 2, False),
            ColumnTestCase(141, ColumnType.TEXT, 64, False),
        ]
        
        for test in tests:
            try:
                column = Column.from_int(test.value)
                self.assertEqual(column.type, test.expected)
                self.assertEqual(column.length, test.expected_length)
                self.assertEqual(test.throws_error, False)
            except Exception as e:
                if not test.throws_error:
                    import pudb; pudb.set_trace()
                    self.fail(e)

    def test_column_int_reading(self):
        for i in range(-1, 127):
            self.assertEqual(Column.from_int(i).to_int(), i)

    def test_read_value(self):
        # read each type of value
        tests = [
            ReadValueTestCase(
                Column.from_int(0),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                1,
                False
            ),
            ReadValueTestCase(
                Column.from_int(1),
                bytes([0x12, 0x34, 0x56]),
                1,
                52, # 0x34
                2,
                False,
            ),
            ReadValueTestCase(
                Column.from_int(2),
                bytes([0x01, 0x02, 0x03, 0x04]),
                1,
                515, # 0x0203
                3,
                False,
            ),
            ReadValueTestCase(
                Column.from_int(3),
                bytes([0x01, 0x02, 0x03, 0x04, 0x05]),
                1,
                131844, # 0x020304
                4,
                False,
            ),
            ReadValueTestCase(
                Column.from_int(4),
                bytes([0x01, 0x02, 0x03, 0x04, 0x05, 0x06]),
                1,
                33752069, # 0x02030405
                5,
                False,
            ),
            ReadValueTestCase(
                Column.from_int(5),
                bytes([
                    0x01, 0x02, 0x03, 0x04,
                    0x05, 0x06, 0x07, 0x08,
                    0x09, 0x10, 0x11, 0x12,
                ]),
                2,
                3315799033608, # 0x030405060708
                8,
                False,
            ),
            ReadValueTestCase(
                Column.from_int(6),
                bytes([
                    0x01, 0x02, 0x03, 0x04,
                    0x05, 0x06, 0x07, 0x08,
                    0x09, 0x10, 0x11, 0x12,
                ]),
                1,
                144964032628459529,
                9,
                False,
            ),
            # IEEE754 int unsupported
            ReadValueTestCase(
                Column.from_int(7),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # constant value 0 
            ReadValueTestCase(
                Column.from_int(8),
                bytes([0x12, 0x34, 0x56]),
                1,
                0,
                1,
                False,
            ),
            # constant value 1  
            ReadValueTestCase(
                Column.from_int(9),
                bytes([0x12, 0x34, 0x56]),
                1,
                1,
                1,
                False,
            ),
            # error if trying to use reserved column type
            ReadValueTestCase(
                Column.from_int(10),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # error if trying to use reserved column type
            ReadValueTestCase(
                Column.from_int(11),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # read string example
            ReadValueTestCase(
                Column.from_int(107),
                bytes([0x12, 0x34, 0x56]) + bytes('it was the faintest 小战俘 one had ever seen', 'utf-8') + bytes([0x11]),
                3,
                'it was the faintest 小战俘 one had ever seen',
                50,
                False,
            ),
            # read blob
            ReadValueTestCase(
                Column.from_int(18),
                bytes([0x12, 0x34, 0x56, 0x78, 0x90]),
                1,
                bytes([0x34, 0x56, 0x78]),
                4,
                False,
            ),

            # unknown column type
            ReadValueTestCase(
                Column.from_int(-1),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
        ]
        for test in tests:
            try:
                value, cursor = Record.read_value(
                    test.column.type,
                    test.data,
                    test.cursor_start,
                    test.column.length,
                )
                self.assertEqual(value, test.expected)
                self.assertEqual(cursor, test.expected_cursor_end)
            except Exception as e:
                if not test.throws_error:
                    self.fail(e)

if __name__ == '__main__':
    unittest.main()
