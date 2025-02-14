export interface User {
	id: string;
	defaultUser: boolean;
	createdAt: string;
	updatedAt: string | null;
}

export interface GetOrCreateDefaultUserResponse {
	user: User;
}
