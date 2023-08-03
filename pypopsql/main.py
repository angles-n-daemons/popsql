from btree import Node
from header import Header
from pager import Pager

def test_btree():
    p = Pager('test.db')
    data = p.get_page(2)
    n = Node(data)
    n._debug_print_cells()

def test_pager():
    p = Pager('test.db')
    stuff = p.get_page(2)

def test_header():
    p = Pager('test.db')
    data = p.get_page(1)
    header = Header(data[:100])
    header._debug_print_values()

def test_schema_btree_page():
    p = Pager('test.db')
    data = p.get_page(1)
    n = Node(data, True)
    n._debug_print_cells()

if __name__ == '__main__':
    test_schema_btree_page()
