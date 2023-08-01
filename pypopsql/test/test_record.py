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

class TestColumnType(TestCase):
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

if __name__ == '__main__':
    unittest.main()
