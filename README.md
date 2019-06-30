# axon
Extremely simple IPC for systems that implement POSIX semaphores and shared memory. 

This library is meant to perform one simple task. Communication between two processes on the same system as quickly and efficiently as possible. It provides an interface to open communication between the two processes and to send and receive byte buffers between them. It allows messages to be tagged with a name, so that the processes can differentiate them from eachother. The goal is for this library to be very small and for its only requirement to be POSIX. Another goal is to write an implementation for Python, C and Go.

## Why?

I have recently run into the need to have fast IPC between slow Python code and a faster language for cpu intensive calculations. My initial use case for this is a Python program that converts [DFT](https://en.wikipedia.org/wiki/Density_functional_theory) data into training sets for neural networks. This relates to the following [research](https://www.nature.com/articles/s41467-019-10343-5.pdf?origin=ppub), for those interested.

## How does it work?

It uses named semaphores to facilitate communication between processes. This allows two processes to synchronize their actions simply by using the same string to create a named semaphore. It then sends the actual data by having the sending process allocate shared memory and trigger the semahpore to inform the receiving process that it has memory available for it to access.

## But Adam, (insert library here) already does this!

I'm sure it does. For better or worse, I like to reinvent the wheel sometimes. Knowing every line of code allows me to make modifications for evantualities that another author might not have predicted. In my experience, such eventualities occur much more often than I would like. I have been forced to write something myself because existing libraries worked at first but later became impractical because of changing requirements. A side goal of this library is for it to be as small as possible so that it is easy to understand and modify.
