# CrowdSort - A Sorting Network Implementation Using Go

![](/docs/images/peek.2.gif)

## Description

This project serves as a proof-of-concept for implementing parallel-sorting algorithms to enhance crowd-sourcing tasks. Specifically, it leverages the theory of [sorting networks](https://en.wikipedia.org/wiki/Sorting_network) to efficiently distribute and manage common crowd-sourcing activities, such as ranking or labeling datasets. By applying sorting network theories, the project aims to demonstrate how parallelism can optimize performance and scalability in crowd-sourcing workflows. Currently, only [Batcher-Even-Odd Merge Sort](https://en.wikipedia.org/wiki/Batcher_odd%E2%80%93even_mergesort) is supported.

## Architecture

![](/docs/images/arch.png)

The implementation uses a simple architecture wth three primary components:

1. `io`: This module serves as an interface between the client and the library. Any form of control and information exchange is directly handled by this module. At its current state, it provides API for standard inputs and outputs.

2. `selector`: This module selects a parallel-sorting algorithm and constructs a dependency graph to create a thread-safe and stateful model for dispatching elements for comparison

3. `dispatcher`: This module efficiently distributes elements for concurrent comparision by units (called Comparators) which may take variable amount of time to yield a result.

## How It All Works

Sorting networks are a fascinating and efficient method for sorting sequences of items using **a fixed sequence of comparisons**. Unlike traditional sorting algorithms (e.g., quicksort or mergesort), which adapt their comparisons based on the input, sorting networks perform the same set of comparisons regardless of the input data. This makes them highly parallelizable and suitable for hardware implementations.

A sorting network consists of two main components:

- _Comparators_: These are entities that take two inputs, compare them, and swap them if they are out of order.

- _Wires_: These carry the values between comparators.

![Source: Wikepedia ](https://upload.wikimedia.org/wikipedia/commons/e/e8/Sorting-network-comparator-demonstration.svg)

A more detailed account on this subject has provided by Knuth in [The Art of Computer Programming Vol. 3](https://github.com/Amanuel94/math-dumps/blob/main/computer-science/The.Art.of.Computer.Programming.3.Sorting.and.Searching.pdf.1), Section 5.3.4.

The overall process can be explained in a few steps.

### Step 1: Registration

A client provides the items to be compared for the `io` module. The `io` module will wrap the given items using the `Wire` to associate them with a wire. `Wire` implements an interface `Comparable` that defines a method `Compare` which the client can define. All comparisions are done via this method. The client also provides the comparator modules. The `io` also wraps these inputs by the `Comparator` interface that defines `CompareEntries` method whose logic can be defined by a client.

### Step 2: Pair Generation

Sorting networks work by comparing wires in pairs using a parallel-sorting. Currently, only BEMS is supported but the implementation was done taking possible future introduction of other algorithms in mind. The `selector` queues the appropriate indices and queues them in data structures called `Connector`s. These data structures will be issued ids along with the corresponding `Wire` pairs. The ids will be used to ease the dependency graph implementation.

### Step 3: Dependency Graph Creation

Eventhough, sorting networks are optimized for parallel sortings tasks, not all comparisions happen at once. Because of that reason `Connector`s have one-way dependency which is modeled using a [DAG](https://stumash.github.io/Algorithm_Notes/trees_graphs/topsort/topsort.html) in this implementation. This will help the dispatcher not worry about false comparision result due to delayed comparator outputs.

### Step 4: Task Assignment

After the graph is created, the dispatcher will select the `Connector`s with 0 dependenceies and assignes them `Comparator`s that are at the moment idle and have capacity to perform more comparisions. This capacity is decided by the client. The `Comparator` are assigned based on their task count (prioritizing those that did less number of tasks) and for that an mutex-protected priority queue is implemented.

## Use Cases

### Data Ranking

CrowdSort can be used to efficiently rank items such as survey responses, products, or user preferences in crowd-sourced datasets. By employing parallel sorting networks, the system can handle large volumes of data and produce rankings quickly, ensuring that decisions or recommendations based on these rankings are delivered in a timely manner. This makes it ideal for applications like market research, product ranking, or social media trend analysis.

### Labeling Tasks

In tasks requiring large-scale dataset labeling, where clear catagorical distiniction between examples is more difficult than relative categorization, CrowdSort can be used to distribute a more robust crowdsorcing labelling.

### Workforce Optimization

CrowdSort can organize and prioritize tasks for crowd-sourced workforces, ensuring efficient assignment and execution. Tasks such as reviewing content, conducting quality checks, or resolving user requests can be ranked and distributed based on urgency or importance, maximizing the productivity of distributed teams while maintaining consistency and fairness in task allocation.

## Example

The example at the beginning (whose implementation is shown [here](/example/main.go)) shows how three comparators are assigned tasks of comparing items for sorting them from smallest to largest cocurrently while providing real-time information on the status of the comparators. The example simulates asynchronous behaviour by pausing the compartors' function for a random amount of time between 1 and 3 seconds. When the status of a wire is `COMPLETED`, it means the corresponding value is at the correct ranking. For demonstration purposes, numerical items are shown in example. However, the implementation was done having any [linearly ordered set](https://en.wikipedia.org/wiki/Total_order) of items in mind.

## Usage

```bash
go get github.com/Amauel94/crowdsort
```
