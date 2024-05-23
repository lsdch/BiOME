import { Edge, Node, useVueFlow } from '@vue-flow/core';
import dagre from 'dagre';
import { ref } from 'vue';
import { ConnectedGroup } from './habitat_graph';

export type NodeData = (ConnectedGroup) & { cluster: string }

/**
 * Composable to run the layout algorithm on the graph.
 * It uses the `dagre` library to calculate the layout of the nodes and edges.
 */
export function useLayout() {

  const blockLayout = ref(false)
  const { findNode } = useVueFlow()

  const graph = ref(new dagre.graphlib.Graph())


  function layout(nodes: Node<NodeData>[], edges: Edge[], direction: 'TB' | 'BT' | 'LR' | 'RL') {
    if (blockLayout.value) {
      blockLayout.value = false
      return nodes
    }
    // we create a new graph instance, in case some nodes/edges were removed, otherwise dagre would act as if they were still there
    const dagreGraph = new dagre.graphlib.Graph({ compound: true })

    graph.value = dagreGraph

    dagreGraph.setDefaultEdgeLabel(() => ({}))

    dagreGraph.setGraph({ rankdir: direction, ranker: 'network-simplex', nodesep: 20, ranksep: 60, edgesep: 10 })


    nodes.forEach((node) => {
      const graphNode = findNode(node.id)

      // Cluster as pseudo node
      dagreGraph.setNode(`cluster.${node.id}`, { label: `cluster.${node.id}`, width: graphNode?.dimensions.width || 150, height: graphNode?.dimensions.height || 50 })
      // Actual group node
      dagreGraph.setNode(node.id, { label: node.id, width: graphNode?.dimensions.width || 150, height: 50 })
      dagreGraph.setParent(node.id, `cluster.${node.id}`)
      node.data?.elements.forEach((child) => {
        dagreGraph.setNode(child.id, { label: child.id, width: graphNode?.dimensions.width || 150, height: 20 })
        dagreGraph.setParent(child.id, `cluster.${node.id}`)
      })
    })

    edges.forEach(edge => dagreGraph.setEdge(edge.sourceHandle!, edge.target))
    dagre.layout(dagreGraph, { compound: true })
    return nodes.map<Node<NodeData>>((node) => {
      const nodeWithPosition = dagreGraph.node(`cluster.${node.id}`)
      node.data?.elements.sort((a, b) => {
        const posA = dagreGraph.node(a.id)
        const posB = dagreGraph.node(b.id)
        return posA.y - posB.y
      })
      return { ...node, position: { x: nodeWithPosition.x, y: nodeWithPosition.y } }
    })
  }

  return {
    /**
     * Graph representation using `dagre`
     */
    graph,
    /**
     * Trigger automatic layout of graph nodes
     */
    layout,
    /**
     * Block next automatic layout attempt
     */
    blockLayout,
  }
}
