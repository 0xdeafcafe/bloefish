import { Button, Icon, Input, Text, VStack } from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';
import { DialogActionTrigger, DialogBody, DialogCloseTrigger, DialogContent, DialogFooter, DialogHeader, DialogRoot, DialogTitle, DialogTrigger } from '~/components/ui/dialog';
import { conversationApi } from '~/api/bloefish/conversation';
import { LuTrash2 } from 'react-icons/lu';
import { InteractionActionButton } from '../atoms/InteractionActionButton';

interface DeleteInteractionDialogProps {
	interactionId: string;
	onDeleteSuccess?: () => void;
}

export const DeleteInteractionDialog: React.FC<DeleteInteractionDialogProps> = ({
	interactionId,
	onDeleteSuccess,
}) => {
	const [isDeleting, setIsDeleting] = useState(false);
	const [confirmText, setConfirmText] = useState('');
	const initialFocusRef = useRef<HTMLInputElement>(null);
	const dialogCloseRef = useRef<HTMLButtonElement>(null);
	const [deleteInteractions] = conversationApi.useDeleteInteractionsMutation();

	const deleteAllowed = confirmText.toLowerCase() === 'delete';

	const handleDelete = async () => {
		if (!deleteAllowed) return;
		if (isDeleting) return;

		setIsDeleting(true);

		try {
			await deleteInteractions({
				interactionIds: [interactionId],
			}).unwrap();

			dialogCloseRef.current?.click();
			onDeleteSuccess?.();
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
				<InteractionActionButton danger tooltip="Delete message">
					<LuTrash2 />
				</InteractionActionButton>
			</DialogTrigger>
			<DialogContent onKeyDown={(e) => e.key === 'Enter' && deleteAllowed && handleDelete()}>
				<DialogHeader>
					<DialogTitle>Delete message</DialogTitle>
				</DialogHeader>
				<DialogBody>
					<VStack align="stretch" gap={6}>
						<Text>
							You're about to delete this message.
							This action cannot be undone and will permanently remove it from the conversation.
						</Text>
						<Input
							ref={initialFocusRef}
							disabled={isDeleting}
							autoFocus
							placeholder="Type 'delete' to confirm"
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
