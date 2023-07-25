from pager import Pager
from btree import Node

def test_btree():
    p = Pager('/tmp/test.db')
    data = p.get_page(2)
    n = Node(data)
    n.test_iter_cell_pointers()


def test_pager():
    p = Pager('/tmp/test.db')
    stuff = p.get_page(2)
    import pudb;
    pudb.set_trace()

if __name__ == '__main__':
    test_btree()
