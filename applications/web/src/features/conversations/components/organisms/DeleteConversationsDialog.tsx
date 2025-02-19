import { Button, Icon, Input, Text, VStack } from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';
import { DialogActionTrigger, DialogBody, DialogCloseTrigger, DialogContent, DialogFooter, DialogHeader, DialogRoot, DialogTitle, DialogTrigger } from '~/components/ui/dialog';
import { conversationApi } from '~/api/bloefish/conversation';
import { LuTrash2 } from 'react-icons/lu';

interface DeleteConversationsDialogProps {
	conversationIds: string[];
	onDeleteSuccess: () => void;
	deleteButtonSize: '2xs' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'; // ButtonVariant -> styled-system
	deleteButtonIconSize: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'; // IconVariant -> styled-system
	deleteButtonText: string;
}

export const DeleteConversationsDialog: React.FC<DeleteConversationsDialogProps> = ({
	conversationIds,
	onDeleteSuccess,
	deleteButtonSize,
	deleteButtonIconSize,
	deleteButtonText,
}) => {
	const [isDeleting, setIsDeleting] = useState(false);
	const [confirmText, setConfirmText] = useState('');
	const initialFocusRef = useRef<HTMLInputElement>(null);
	const dialogCloseRef = useRef<HTMLButtonElement>(null);
	const [deleteConversations] = conversationApi.useDeleteConversationsMutation();

	const deleteAllowed = confirmText.toLowerCase() === 'delete';

	const handleDelete = async () => {
		if (!deleteAllowed) return;
		if (isDeleting) return;

		setIsDeleting(true);

		try {
			await deleteConversations({
				conversationIds,
			}).unwrap();

			dialogCloseRef.current?.click();
			onDeleteSuccess();
		} catch (e) {
			console.error(e);
		} finally {
			setIsDeleting(false);
			setConfirmText('');
		}
	};

	useEffect(() => {
		initialFocusRef.current?.focus();
	}, []);

	return (
		<DialogRoot role="alertdialog">
			<DialogTrigger asChild>
				<Button
					variant={'outline'}
					colorPalette={'red'}
					size={deleteButtonSize}
				>
					<Icon size={deleteButtonIconSize}>
						<LuTrash2 />
					</Icon>
					{deleteButtonText}
				</Button>
			</DialogTrigger>
			<DialogContent onKeyDown={(e) => e.key === 'Enter' && deleteAllowed && handleDelete()}>
				<DialogHeader>
					<DialogTitle>
						{`Delete conversation${conversationIds.length > 1 ? 's' : ''}`}
					</DialogTitle>
				</DialogHeader>
				<DialogBody>
					<VStack align="stretch" gap={6}>
						<Text>
							You're about to delete {conversationIds.length} conversation{conversationIds.length === 1 ? '' : 's'}.
							This action cannot be undone and will permanently remove all selected conversations from your account.
						</Text>
						<Input
							ref={initialFocusRef}
							disabled={isDeleting}
							autoFocus
							placeholder={'Type \'delete\' to confirm'}
							value={confirmText}
							onChange={(e) => setConfirmText(e.target.value)}
							autoComplete="off"
						/>
					</VStack>
				</DialogBody>
				<DialogFooter>
					<DialogActionTrigger asChild>
						<Button variant="outline">Cancel</Button>
					</DialogActionTrigger>
					<Button
						colorPalette="red"
						disabled={!deleteAllowed || isDeleting}
						loading={isDeleting}
						onClick={handleDelete}
					>
						Delete
					</Button>
				</DialogFooter>
				<DialogCloseTrigger ref={dialogCloseRef} />
			</DialogContent>
		</DialogRoot>
	);
};
