import { Button, Input, Text, VStack } from '@chakra-ui/react';
import { useRef, useState } from 'react';
import { DialogActionTrigger, DialogBody, DialogCloseTrigger, DialogContent, DialogFooter, DialogHeader, DialogRoot, DialogTitle, DialogTrigger } from '~/components/ui/dialog';
import { conversationApi } from '~/api/bloefish/conversation';
import { LuTrash2 } from 'react-icons/lu';
import { InteractionActionButton } from '../atoms/InteractionActionButton';

interface DeleteInteractionDialogProps {
	disabled?: boolean,
	interactionId: string;
	onDeleteSuccess?: () => void;
}

export const DeleteInteractionDialog: React.FC<DeleteInteractionDialogProps> = ({
	disabled,
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

	return (
		<DialogRoot role={'alertdialog'} initialFocusEl={() => initialFocusRef.current}>
			<DialogTrigger asChild>
				<InteractionActionButton
					disabled={disabled}
					danger
					tooltip={'Delete message'}
				>
					<LuTrash2 />
				</InteractionActionButton>
			</DialogTrigger>
			<DialogContent onKeyDown={(e) => e.key === 'Enter' && deleteAllowed && handleDelete()}>
				<DialogHeader>
					<DialogTitle>{'Delete message'}</DialogTitle>
				</DialogHeader>
				<DialogBody>
					<Text mb={6}>
						{'You\'re about to delete a message. This action cannot be '}
						{'undone and will permanently remove the selected message from '}
						{'from the platform.'}
					</Text>
					<VStack align={'stretch'} gap={2}>
						<Text>
							{'Are you sure you with to proceed?'}
						</Text>
						<Input
							ref={initialFocusRef}
							disabled={isDeleting}
							autoFocus
							placeholder={'Type \'delete\' to confirm'}
							value={confirmText}
							onChange={(e) => setConfirmText(e.target.value)}
							autoComplete={'off'}
						/>
					</VStack>
				</DialogBody>
				<DialogFooter>
					<DialogActionTrigger asChild>
						<Button variant={'outline'}>Cancel</Button>
					</DialogActionTrigger>
					<Button
						colorPalette={'red'}
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
