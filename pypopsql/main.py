from pager import Pager
from btree import Node

def test_btree():
    p = Pager('test.db')
    data = p.get_page(2)
    n = Node(data)
    n._debug_print_cells()


def test_pager():
    p = Pager('test.db')
    stuff = p.get_page(2)

if __name__ == '__main__':
    test_btree()
