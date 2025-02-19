import { Icon, IconButton } from "@chakra-ui/react";
import type React from "react";
import { Tooltip } from "~/components/ui/tooltip";

interface InteractionActionButtonProps extends React.PropsWithChildren {
	loading?: boolean;
	disabled?: boolean;
	tooltip?: string;
	danger?: boolean;
	onClick?: () => void;
}

export const InteractionActionButton: React.FC<InteractionActionButtonProps> = ({
	loading,
	children,
	disabled,
	tooltip,
	danger,
	onClick,
}) => {
	return (
		<Tooltip content={tooltip}>
			<IconButton
				colorPalette={danger ? 'red' : void 0}
				variant={'ghost'}
				size={'2xs'}
				disabled={disabled}
				loading={loading}
				onClick={onClick}
			>
				<Icon size={'xs'}>
					{children}
				</Icon>
			</IconButton>
		</Tooltip>
	);
};
