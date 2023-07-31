from unittest import TestCase
from dataclasses import dataclass
from util import varint


@dataclass
class VarIntTestCase:
    data: bytes
    expected: int
    num_bytes: int


class TestVarint(TestCase):
    def test_varint(self):
        tests = [
            VarIntTestCase([0x00], 0x00000000, 1),
            VarIntTestCase([0x7f], 0x0000007f, 1),
            VarIntTestCase([0x81, 0x00], 0x00000080, 2),
            VarIntTestCase([0x80, 0x7f], 0x0000007f, 2),
            VarIntTestCase([0x81, 0x91, 0xd1, 0xac, 0x78], 0x12345678, 5),
            VarIntTestCase([0x81, 0x81, 0x81, 0x81, 0x01], 0x10204081, 5),
        ]
        for test in tests:
            with self.subTest(msg=f'testing varint {test.data} = {test.expected}'):
                result, num_bytes = varint(test.data, 0)
                self.assertEqual(test.expected, result)
                self.assertEqual(test.num_bytes, num_bytes)

    def test_varint_mid_sequence(self):
        result, num_bytes = varint([
            0x32, 0x80, 0x7f, 0x91,
        ], 1)
        self.assertEqual(result, 0x0000007f)
        self.assertEqual(num_bytes, 2)

    def test_varint_full_9_bytes(self):
        result, num_bytes = varint([
            0x81, 0x81, 0x81, 0x81, 0x81, 0x81, 0x81, 0x81, 0x81
        ], 0)
        self.assertEqual(result, 145249953336295809)
        self.assertEqual(num_bytes, 9)

if __name__ == '__main__':
    unittest.main()
