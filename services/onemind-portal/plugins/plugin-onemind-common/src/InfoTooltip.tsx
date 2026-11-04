import { ReactNode } from 'react';
import { IconButton, Tooltip } from '@material-ui/core';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';

export type InfoTooltipProps = {
  /** Name of the card, for the button label. */
  label: string;
  /** Where the data comes from and when it was observed, plus what it does not prove. */
  children: ReactNode;
};

/**
 * One per card. Replaces the repeated "this is not proof of health" paragraphs: the caveat stays
 * available (hover, focus or tap) without taking card space.
 */
export function InfoTooltip({ label, children }: InfoTooltipProps) {
  return (
    <Tooltip title={<>{children}</>} interactive enterTouchDelay={0} leaveTouchDelay={8000}>
      <IconButton size="small" aria-label={`About ${label}`} style={{ padding: 4 }}>
        <InfoOutlinedIcon fontSize="small" />
      </IconButton>
    </Tooltip>
  );
}

/** Compact "source · observed time" line shown once per card. */
export function SourceLine({ source, observedAt }: { source: string; observedAt?: string }) {
  return (
    <span style={{ fontSize: 12, opacity: 0.8 }}>
      {source}
      {observedAt ? ` · observed ${new Date(observedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}` : ''}
    </span>
  );
}
