# Virtual Machine

The virtual machine is a fairly straightforward implementation. It consists of the following properties:

 - *registers* a list of virtual registers who can hold any type of value.
 - *counter* the counter to the current instruction.

The virtual machine is *not* thread safe and therefore should not be accessed by multiple goroutines at the same time.

To run the virtual machine, it needs to be given a program. A program in this example is just a list of [instructions](documentation/core/instruction.md).

Once running the VM it will continue to execute instructions until it reaches the halt operator. On each instruction it executes it will increment its internal counter by one - unless the instruction explicitly requests a next instruction for execution.
