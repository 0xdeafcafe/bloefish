import type { Actor } from './shared.types';

export interface SkillSet {
	id: string;
	name: string;
	description: string;
	icon: string;
	prompt: string;

	owner: Actor;

	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
}

export interface CreateSkillSetRequest {
	name: string;
	icon: string;
	description: string;
	prompt: string;
	owner: Actor;
}

export interface GetSkillSetRequest {
	id: string;
	name: string;
	icon: string;
	description: string;
	prompt: string;

	owner: Actor;

	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
}

export type GetSkillSetResponse = SkillSet;

export interface ListSkillSetsByOwnerRequest {
	owner: Actor;
}

export interface ListSkillSetsByOwnerResponse {
	skillSets: SkillSet[];
}
