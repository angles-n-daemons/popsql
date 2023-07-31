from pager import Pager
from btree import Node

def test_btree():
    p = Pager('test.db')
    data = p.get_page(2)
    n = Node(data)
    cells = n.cells
    for cell in cell:
        print('row id: ', cell.row_id)
        print('payload size: ', cell.payload_size)
        print('payload: ', cell.payload)
        print('cursor end: ', cell.cursor)
        print('\n\n')


def test_pager():
    p = Pager('test.db')
    stuff = p.get_page(2)

if __name__ == '__main__':
    test_btree()
