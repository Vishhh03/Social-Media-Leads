import React from 'react';
import { Handle, Position } from '@xyflow/react';

export function AINode({ data }) {
    return (
        <div className="node-card" style={{ borderColor: '#d946ef' }}>
            <Handle
                type="target"
                position={Position.Top}
                className="node-handle node-handle-ai"
            />
            <div className="node-header node-header-ai fuchsia">
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2.5"><path d="M12 2a10 10 0 1 0 10 10H12V2z" /><path d="M12 12L2.05 9.27" /><path d="M12 12l7.07 7.07" /></svg>
                    AI Agent
                </div>
                <span style={{ background: 'rgba(255,255,255,0.2)', padding: '2px 6px', borderRadius: '4px', fontSize: '10px' }}>GPT-4o</span>
            </div>
            <div className="node-body">
                <div className="node-title">
                    {data.label || 'Generate Reply'}
                </div>
                <div className="node-desc" style={{ marginBottom: '12px', display: 'block' }}>
                    {data.description || 'Reads Knowledge Base & replies automatically'}
                </div>
                <div className="node-prompt">
                    {data.prompt || 'Prompt: You are a helpful assistant...'}
                </div>
            </div>

            {/* If this AI node routes based on intent, it will have multiple outputs */}
            {data.isRouter ? (
                <div style={{ display: 'flex', justifyContent: 'space-evenly', paddingBottom: '8px', position: 'relative', height: '16px' }}>
                    <Handle type="source" position={Position.Bottom} id="hot" className="node-handle node-handle-ai" style={{ position: 'sticky', left: '25%' }} />
                    <Handle type="source" position={Position.Bottom} id="cold" className="node-handle node-handle-ai" style={{ position: 'sticky', left: '75%' }} />
                </div>
            ) : (
                <Handle type="source" position={Position.Bottom} className="node-handle node-handle-ai" />
            )}
        </div>
    );
}
