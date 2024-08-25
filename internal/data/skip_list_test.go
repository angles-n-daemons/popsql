package data_test

// create a benchmark simulating likely usage

// test using a single head which is always the height of the tree
// -- report how many nodes of each height there are
// -- test alternative where there are multiple heads
// -- check insert performance
// -- check lookup performance

// test pointer vs value performance (mem, speed)
// - with a linked list, just creating a lot of nodes
//    - create a lot of goroutines, search through all of them
// - node.next, prev
// - list.head(s)

// test performance using generics vs not generics

// add unit tests

// question around heights
// oh if I use the approach where I always insert new head nodes if the node is smaller, does it drastically increase the size of the lists, number of nodes per level?
// I can use a random seed to figure this out

// test skiplist byte key, byte string

// what if I use int8 for height (performance for code writability)
