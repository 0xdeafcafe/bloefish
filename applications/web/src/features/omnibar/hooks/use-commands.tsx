import { useNavigate } from 'react-router';
import { LuFolderOpen, LuMessageCirclePlus } from 'react-icons/lu';

export interface Command {
	name: string;
	keywords?: string[];
	command?: string;
	onInvoke: () => void;
	icon: React.ReactElement,
}

export function useCommands(): Command[] {
	const navigate = useNavigate();

	return [{
		name: 'New conversation...',
		onInvoke: () => navigate('/'),
		icon: <LuMessageCirclePlus />,
	}, {
		name: 'View conversations',
		onInvoke: () => navigate('/'),
		icon: <LuFolderOpen />,
	}];
}
