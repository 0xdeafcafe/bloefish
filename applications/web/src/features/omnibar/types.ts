import type { Interaction } from '../conversations/store/types';
import type { Command } from './hooks/use-commands';

export type SearchContextType = 'command' | 'interaction';

export interface CommandSearchContext extends Command {
	searchContextType: 'command';
}

export interface InteractionSearchContext extends Interaction {
	searchContextType: 'interaction';
}

export type SearchContextItem = CommandSearchContext | InteractionSearchContext;
