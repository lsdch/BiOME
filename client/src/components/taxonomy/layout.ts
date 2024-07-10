import { Edge, Node, useVueFlow } from '@vue-flow/core';
import dagre from '@dagrejs/dagre';
import { ref } from 'vue';


/**
 * Composable to run the layout algorithm on the graph.
 * It uses the `dagre` library to calculate the layout of the nodes and edges.
 */
// import { ELK, ElkNode } from "elkjs/lib/elk-api"
import ELK, { ElkNode } from 'elkjs'
const elk = new ELK({
  defaultLayoutOptions: {
    'elk.algorithm': 'layered',
    'elk.spacing.nodeNode': '5',
    'elk.layered.spacing.nodeNodeBetweenLayers': '5',
    'elk.separateConnectedComponents': 'false'
  }
})

export function useLayout<NodeData>() {

  const { findNode } = useVueFlow()




  const graph = ref(new dagre.graphlib.Graph())


  async function layout(nodes: Node<NodeData>[], edges: Edge[], direction: 'TB' | 'BT' | 'LR' | 'RL'): Promise<Node<NodeData>[]> {


    const g: ElkNode = {
      id: "root",
      children: nodes.map(({ id }) => {
        return { id, width: 300, height: 30 }
      }),
      edges: edges.map(({ id, source, target }) => {
        return { id, sources: [source], targets: [target] }
      })
    }

    const l = await elk.layout(g)

    return l.children!.map<Node<NodeData>>(({ id, x, y }) => {
      return { id, position: { x: x!, y: y! } }
    })

    // l.then(console.log)
    // // we create a new graph instance, in case some nodes/edges were removed, otherwise dagre would act as if they were still there
    // const dagreGraph = new dagre.graphlib.Graph({ compound: false })

    // graph.value = dagreGraph

    // dagreGraph.setDefaultEdgeLabel(() => ({}))

    // dagreGraph.setGraph({ rankdir: direction, ranker: 'tight-tree', ranksep: 1000, edgesep: 50, nodesep: 10 })


    // nodes.forEach((node) => {
    //   const graphNode = findNode(node.id)
    //   // Actual group node
    //   dagreGraph.setNode(node.id, { label: node.id, width: 100 })
    // })

    // edges.forEach(edge => dagreGraph.setEdge(edge.source, edge.target))

    // dagre.layout(dagreGraph, { compound: false })

    // return nodes.map<Node<NodeData>>((node) => {
    //   const nodeWithPosition = dagreGraph.node(node.id)
    //   return { ...node, position: { x: nodeWithPosition.x, y: nodeWithPosition.y } }
    // })
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
  }
}
