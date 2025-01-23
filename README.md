# CrowdSort - Sorting Network Implementation using Go

![](/docs/images/peek.2.gif)

## Description

This project serves as a proof-of-concept for implementing parallel-sorting algorithms to enhance crowd-sourcing tasks. Specifically, it leverages the theory of [sorting networks](https://en.wikipedia.org/wiki/Sorting_network) to efficiently distribute and manage common crowd-sourcing activities, such as ranking or labeling datasets. By applying sorting network theories, the project aims to demonstrate how parallelism can optimize performance and scalability in crowd-sourcing workflows. Currently, only [Batcher-Even-Odd Merge Sort](https://en.wikipedia.org/wiki/Batcher_odd%E2%80%93even_mergesort) is supported.

## Architecture

![](/docs/images/arch.png)

The implementation uses a simple architecture wth three primary components:

1. `io`: This module serves as an interface between the client and the library. Any form of control and information exchange is directly handled by this module. At its current state, it provides API for standard inputs and outputs.

2. `selector`: This module selects a parallel-sorting algorithm and constructs a dependency graph to create a thread-safe and stateful model for selecting elements for comparison.

3. `dispatcher`: This module efficiently distributes elements for concurrent comparision by units (called Comparators) which may take variable amount of time to yield a result

## Use Cases

### Data Ranking

CrowdSort can be used to efficiently rank items such as survey responses, products, or user preferences in crowd-sourced datasets. By employing parallel sorting networks, the system can handle large volumes of data and produce rankings quickly, ensuring that decisions or recommendations based on these rankings are delivered in a timely manner. This makes it ideal for applications like market research, product ranking, or social media trend analysis.

### Labeling Tasks

In tasks requiring large-scale dataset labeling, where clear catagorical distiniction between examples is more difficult than relative categorization, CrowdSort can be used to distribute a more robust crowdsorcing labelling.

### Workforce Optimization

CrowdSort can organize and prioritize tasks for crowd-sourced workforces, ensuring efficient assignment and execution. Tasks such as reviewing content, conducting quality checks, or resolving user requests can be ranked and distributed based on urgency or importance, maximizing the productivity of distributed teams while maintaining consistency and fairness in task allocation.

## Example

The example at the beginning shows how three comparators are assigned tasks of comparing items for sorting them from smallest to largest cocurrently while providing real-time information on the status of the comparators. When the status of a wire is COMPLETED, it means the corresponding value is at the correct ranking. For demonstration purposes, numerical items are shown in example. However, the implementation was done having any [linearly ordered set](https://en.wikipedia.org/wiki/Total_order) of items in mind.

## Usage

```bash
go get github.com/Amauel94/crowdsort
```
