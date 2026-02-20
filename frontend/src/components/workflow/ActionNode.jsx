import React from 'react';
import { Handle, Position } from '@xyflow/react';

export function ActionNode({ data }) {
    return (
        <div className="node-card">
            <Handle
                type="target"
                position={Position.Top}
                className="node-handle node-handle-action"
            />
            <div className="node-header node-header-action">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2.5"><path d="M22 2L11 13" /><path d="M22 2L15 22l-4-9-9-4z" /></svg>
                Action
            </div>
            <div className="node-body">
                <div className="node-title">
                    {data.label || 'Send Message'}
                </div>
                <div className="node-desc">
                    {data.description || 'Sends a simple text reply back to the user'}
                </div>
            </div>
            <Handle
                type="source"
                position={Position.Bottom}
                className="node-handle node-handle-action"
            />
        </div>
    );
}
