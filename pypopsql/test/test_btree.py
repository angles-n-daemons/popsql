from btree import Node
from unittest import TestCase

class TestBTree(TestCase):
    def test_table_leaf(self):
        data = bytes([
            0x0d, # Node Type 13
            0x00, 0x00, # first_freeblock
            0x00, 0x02, # num cells 2
            0x00, 0x11, # cell offset 17
            0x00, # num fragmented bytes

            ## cell pointer array
            # row 1 offset 25
            0x00, 0x19,
            # row 2 offset 17
            0x00, 0x11,

            # null bytes in the center of the page
            0x00, 0x00, 0x00, 0x00, 0x00,

            # row 2
            0x06, # payload size
            0x02, # row id
            0x03, 0x11, 0x01, 0x79, 0x6f, 0x02, # payload

            # row 1
            0x05, # payload size
            0x01, # row id
            0x03, 0x11, 0x09, 0x68, 0x69, # payload
        ])

        node = Node(data)
        self.assertEqual(node.data, data)
        self.assertEqual(node.page_size, 32)
        self.assertEqual(node.num_cells, 2)
        self.assertEqual(node.cell_offset, 17)
        self.assertEquals(node.first_freeblock, 0)
        self.assertEquals(node.num_fragmented_bytes, 0)
        self.assertEquals(node.right_pointer, None)
        self.assertEquals(len(node.cells), 2)

        # check first row
        self.assertEquals(node.cells[0].payload_size, 5)
        self.assertEquals(node.cells[0].row_id, 1)
        self.assertEquals(node.cells[0].payload, bytes([0x03, 0x11, 0x09, 0x68, 0x69]))
        self.assertEquals(node.cells[0].cursor, 32)

        # check second row
        self.assertEquals(node.cells[1].payload_size, 6)
        self.assertEquals(node.cells[1].row_id, 2)
        self.assertEquals(node.cells[1].payload, bytes([0x03, 0x11, 0x01, 0x79, 0x6f, 0x02]))
        self.assertEquals(node.cells[1].cursor, 25)

if __name__ == '__main__':
    unittest.main()
