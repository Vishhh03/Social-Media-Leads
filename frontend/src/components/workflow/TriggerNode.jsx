import React from 'react';
import { Handle, Position } from '@xyflow/react';

export function TriggerNode({ data }) {
    return (
        <div className="node-card">
            <div className="node-header node-header-trigger">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2.5"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" /></svg>
                Trigger
            </div>
            <div className="node-body">
                <div className="node-title">
                    {data.label || 'Incoming Message'}
                </div>
                <div className="node-desc">
                    {data.description || 'Fires when a new DM is received'}
                </div>
            </div>
            <Handle
                type="source"
                position={Position.Bottom}
                className="node-handle node-handle-trigger"
            />
        </div>
    );
}
