import { Button, ButtonGroup, Input, Kbd, Text, VStack } from '@chakra-ui/react';
import { useRef, useState } from 'react';
import { DialogActionTrigger, DialogBody, DialogCloseTrigger, DialogContent, DialogFooter, DialogHeader, DialogRoot, DialogTitle, DialogTrigger } from '~/components/ui/dialog';

interface DeleteConversationsDialogProps {
  conversationIds: string[];
}

export const DeleteConversationsDialog: React.FC<DeleteConversationsDialogProps> = ({
  conversationIds,
}) => {
  const [isDeleting, setIsDeleting] = useState(false);
  const [confirmText, setConfirmText] = useState('');
  const initialFocusRef = useRef(null);

  const handleDelete = async () => {
    setIsDeleting(true);
    setIsDeleting(false);
  };

  const isDeleteEnabled = confirmText.toLowerCase() === 'delete';

  return (
    <DialogRoot role="alertdialog">
      <DialogTrigger asChild>
        <Button 
          variant={'outline'} 
          colorPalette={'red'} 
          size={'xs'}
        >
          Delete conversation{conversationIds.length > 1 ? 's' : ''}
        </Button>
      </DialogTrigger>
      <DialogContent onKeyDown={(e) => e.key === 'Enter' && isDeleteEnabled && handleDelete()}>
        <DialogHeader>
          <DialogTitle>Delete Conversations</DialogTitle>
        </DialogHeader>
        <DialogBody>
          <VStack align="stretch">
            <Text>
              You're about to delete {conversationIds.length} conversation{conversationIds.length === 1 ? '' : 's'}.
              This action cannot be undone and will permanently remove all selected conversations from your account.
            </Text>
            <Input
              ref={initialFocusRef}
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
            disabled={!isDeleteEnabled}
            loading={isDeleting}
            onClick={handleDelete}
          >
            Delete
          </Button>
        </DialogFooter>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  );
};
