from unittest import TestCase

from btree import Node
from pager import Pager
from record import Column, ColumnType

class TestBTree(TestCase):
    def test_table_leaf_parse(self):
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
        self.assertEqual(node.first_freeblock, 0)
        self.assertEqual(node.num_fragmented_bytes, 0)
        self.assertEqual(node.right_pointer, None)
        self.assertEqual(len(node.cells), 2)

        # check first row
        self.assertEqual(node.cells[0].payload_size, 5)
        self.assertEqual(node.cells[0].row_id, 1)
        self.assertEqual(node.cells[0].payload, bytes([0x03, 0x11, 0x09, 0x68, 0x69]))
        self.assertEqual(node.cells[0].cursor, 32)

        # check second row
        self.assertEqual(node.cells[1].payload_size, 6)
        self.assertEqual(node.cells[1].row_id, 2)
        self.assertEqual(node.cells[1].payload, bytes([0x03, 0x11, 0x01, 0x79, 0x6f, 0x02]))
        self.assertEqual(node.cells[1].cursor, 25)

    def test_schema_header_page(self):
        # test using example database written by sqlite
        # to recreate use the following commands
        # $ sqlite3 test/test.db
        # > CREATE TABLE test(col1 VARCHAR(2), col2 INT);

        p = Pager('./test/test.db')
        data = p.get_page(1)
        node = Node(data, True)
        node._debug_print_cells()
        self.assertEqual(node.page_size, 4096)
        self.assertEqual(node.num_cells, 1)
        self.assertEqual(node.cell_offset, 4018)
        self.assertEqual(node.first_freeblock, 0)
        self.assertEqual(node.num_fragmented_bytes, 0)
        self.assertEqual(node.right_pointer, None)
        self.assertEqual(len(node.cells), 1)

        cell = node.cells[0]
        self.assertEqual(cell.payload_size, 64)
        self.assertEqual(cell.row_id, 1)
        self.assertEqual(cell.cursor, 4084)
        self.assertEquals(cell.record.columns, [
            Column(ColumnType.TEXT, 5),
            Column(ColumnType.TEXT, 4),
            Column(ColumnType.TEXT, 4),
            Column(ColumnType.TINYINT),
            Column(ColumnType.TEXT, 44),
        ])
        self.assertEquals(cell.record.values, ['table', 'test', 'test', 2, 'CREATE TABLE test(col1 varchar(2), col2 int)'])
        self.assertEquals(cell.record.cursor, 4084)

if __name__ == '__main__':
    unittest.main()
