from dataclasses import dataclass
from typing import Union
from unittest import TestCase

from record import Record, ColumnType

@dataclass
class ColumnTypeTestCase:
    value: int

    expected: ColumnType
    expected_length: Union[int, None]
    throws_error: bool

@dataclass
class ReadValueTestCase:
    column_type: ColumnType
    data: bytes
    cursor_start: int

    expected: any
    expected_cursor_end: int
    throws_error: bool

class TestRecord(TestCase):
    def test_column_type_values(self):
        tests = [
            # test 1-11
            ColumnTypeTestCase(0, ColumnType.NULL, None, False),
            ColumnTypeTestCase(1, ColumnType.TINYINT, None, False),
            ColumnTypeTestCase(2, ColumnType.SMALLINT, None, False),
            ColumnTypeTestCase(3, ColumnType.SMALLISHINT, None, False),
            ColumnTypeTestCase(4, ColumnType.INTEGER, None, False),
            ColumnTypeTestCase(5, ColumnType.BIGGISHINT, None, False),
            ColumnTypeTestCase(6, ColumnType.LONG, None, False),
            ColumnTypeTestCase(7, ColumnType.IEEE754INT, None, False),
            ColumnTypeTestCase(8, ColumnType.ZERO, None, False),
            ColumnTypeTestCase(9, ColumnType.ONE, None, False),
            ColumnTypeTestCase(10, ColumnType.RESERVED_1, None, False),
            ColumnTypeTestCase(11, ColumnType.RESERVED_2, None, False),

            # test less than 1
            ColumnTypeTestCase(-1, None, None, True),
            ColumnTypeTestCase(-10, None, None, True),

            # test blob w multiple values
            ColumnTypeTestCase(12, ColumnType.BLOB, 0, False),
            ColumnTypeTestCase(14, ColumnType.BLOB, 1, False),
            ColumnTypeTestCase(16, ColumnType.BLOB, 2, False),
            ColumnTypeTestCase(140, ColumnType.BLOB, 64, False),

            # test text with multiple values
            ColumnTypeTestCase(13, ColumnType.TEXT, 0, False),
            ColumnTypeTestCase(15, ColumnType.TEXT, 1, False),
            ColumnTypeTestCase(17, ColumnType.TEXT, 2, False),
            ColumnTypeTestCase(141, ColumnType.TEXT, 64, False),
        ]
        
        for test in tests:
            try:
                column_type = ColumnType.from_varint(test.value)
                self.assertEqual(column_type, test.expected)
                self.assertEqual(column_type.length, test.expected_length)
                self.assertEqual(test.throws_error, False)
            except Exception as e:
                if not test.throws_error:
                    self.fail(e)

    def test_read_value(self):
        # read each type of value
        tests = [
            ReadValueTestCase(
                ColumnType.from_varint(0),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                1,
                False
            ),
            ReadValueTestCase(
                ColumnType.from_varint(1),
                bytes([0x12, 0x34, 0x56]),
                1,
                52, # 0x34
                2,
                False,
            ),
            ReadValueTestCase(
                ColumnType.from_varint(2),
                bytes([0x01, 0x02, 0x03, 0x04]),
                1,
                515, # 0x0203
                3,
                False,
            ),
            ReadValueTestCase(
                ColumnType.from_varint(3),
                bytes([0x01, 0x02, 0x03, 0x04, 0x05]),
                1,
                131844, # 0x020304
                4,
                False,
            ),
            ReadValueTestCase(
                ColumnType.from_varint(4),
                bytes([0x01, 0x02, 0x03, 0x04, 0x05, 0x06]),
                1,
                33752069, # 0x02030405
                5,
                False,
            ),
            ReadValueTestCase(
                ColumnType.from_varint(5),
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
                ColumnType.from_varint(6),
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
                ColumnType.from_varint(7),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # constant value 0 
            ReadValueTestCase(
                ColumnType.from_varint(8),
                bytes([0x12, 0x34, 0x56]),
                1,
                0,
                1,
                False,
            ),
            # constant value 1  
            ReadValueTestCase(
                ColumnType.from_varint(9),
                bytes([0x12, 0x34, 0x56]),
                1,
                1,
                1,
                False,
            ),
            # error if trying to use reserved column type
            ReadValueTestCase(
                ColumnType.from_varint(10),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # error if trying to use reserved column type
            ReadValueTestCase(
                ColumnType.from_varint(11),
                bytes([0x12, 0x34, 0x56]),
                1,
                None,
                None,
                True,
            ),
            # read string example
            ReadValueTestCase(
                ColumnType.from_varint(107),
                bytes([0x12, 0x34, 0x56]) + bytes('it was the faintest 小战俘 one had ever seen', 'utf-8') + bytes([0x11]),
                3,
                'it was the faintest 小战俘 one had ever seen',
                50,
                False,
            ),
            # read blob
            ReadValueTestCase(
                ColumnType.from_varint(18),
                bytes([0x12, 0x34, 0x56, 0x78, 0x90]),
                1,
                bytes([0x34, 0x56, 0x78]),
                4,
                False,
            ),

            # unknown column type
            ReadValueTestCase(
                ColumnType.from_varint(-1),
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
                    test.column_type,
                    test.data,
                    test.cursor_start,
                )
                self.assertEqual(value, test.expected)
                self.assertEqual(cursor, test.expected_cursor_end)
            except Exception as e:
                if not test.throws_error:
                    self.fail(e)
        pass

    def test_read_values(self):
        pass
        # read value with no cursor movement
        # read value with string
        # read value with integer

if __name__ == '__main__':
    unittest.main()
