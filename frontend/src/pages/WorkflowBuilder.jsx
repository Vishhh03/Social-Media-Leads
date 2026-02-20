import React, { useState, useCallback } from 'react';
import {
    ReactFlow,
    MiniMap,
    Controls,
    Background,
    useNodesState,
    useEdgesState,
    addEdge,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import '@xyflow/react/dist/style.css';
import { TriggerNode } from '../components/workflow/TriggerNode';
import { ActionNode } from '../components/workflow/ActionNode';
import { AINode } from '../components/workflow/AINode';
import { generateWorkflow } from '../api';

const nodeTypes = {
    trigger_meta_dm: TriggerNode,
    action_send_message: ActionNode,
    action_ai_reply: AINode,
};

const initialNodes = [
    { id: '1', type: 'trigger_meta_dm', position: { x: 250, y: 50 }, data: { label: 'Incoming Meta DM' } },
    { id: '2', type: 'action_ai_reply', position: { x: 250, y: 250 }, data: { label: 'Generate Quote', description: 'Reads pricing sheet and replies', prompt: 'Analyze lead and reply with quote' } },
    { id: '3', type: 'action_send_message', position: { x: 250, y: 450 }, data: { label: 'Follow Up (1 Day)', description: 'Send delayed follow up message' } },
];
const initialEdges = [
    { id: 'e1-2', source: '1', target: '2', animated: true, style: { stroke: '#6366f1' } },
    { id: 'e2-3', source: '2', target: '3', animated: true, style: { stroke: '#6366f1' } },
];

export default function WorkflowBuilder() {
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
    const [prompt, setPrompt] = useState("");

    const [isGenerating, setIsGenerating] = useState(false);

    const onConnect = useCallback(
        (params) => setEdges((eds) => addEdge(params, eds)),
        [setEdges],
    );

    const handleMagicPrompt = async () => {
        if (!prompt) return;
        setIsGenerating(true);
        try {
            const data = await generateWorkflow(prompt);
            setNodes(data.nodes || []);
            setEdges(data.edges || []);
            setPrompt("");
        } catch (error) {
            console.error("Failed to generate workflow:", error);
            alert("Error generating workflow: " + error.message);
        } finally {
            setIsGenerating(false);
        }
    }

    console.log("WorkflowBuilder Rendering:");
    console.log("Nodes ->", nodes);
    console.log("Edges ->", edges);
    console.log("NodeTypes ->", nodeTypes);

    try {
        return (
            <div className="workflow-container" style={{ minHeight: 'calc(100vh - 60px)' }}>
                <div className="workflow-header">
                    <h1 className="workflow-header-title">Workflow Builder (AI Co-Pilot)</h1>
                    <div className="workflow-toolbar">
                        <input
                            type="text"
                            placeholder="E.g. Reply with pricing, wait 1 hr..."
                            value={prompt}
                            onChange={(e) => setPrompt(e.target.value)}
                            className="input"
                        />
                        <button
                            onClick={handleMagicPrompt}
                            disabled={isGenerating}
                            className="btn btn-primary">
                            {isGenerating ? "✨ Generating..." : "✨ Generate"}
                        </button>
                    </div>
                    <button className="btn btn-primary" style={{ backgroundColor: 'var(--success)', borderColor: 'var(--success)' }}>
                        Save Workflow
                    </button>
                </div>

                <div className="workflow-canvas" style={{ minHeight: '600px', height: '100%', border: '2px dashed red' }}>
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        nodeTypes={nodeTypes}
                        onNodesChange={onNodesChange}
                        onEdgesChange={onEdgesChange}
                        onConnect={onConnect}
                        fitView
                        proOptions={{ hideAttribution: true }}
                    >
                        <Controls />
                        <MiniMap
                            nodeStrokeColor={() => '#6366f1'}
                            nodeColor={() => '#e0e7ff'}
                        />
                        <Background variant="dots" gap={12} size={1} />
                    </ReactFlow>
                </div>
            </div>
        );
    } catch (e) {
        return <div style={{ padding: 50, color: 'red' }}><h1>Canvas Render Error:</h1><pre>{e.toString()}</pre></div>;
    }
}
