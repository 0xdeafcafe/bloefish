import { useRef } from 'react';
import { LuChevronDown, LuFilePlus, LuFileUp, LuTrash } from 'react-icons/lu';
import { Button } from '~/components/ui/button';
import { userApi } from '~/api/bloefish/user';
import { PopoverArrow, PopoverBody, PopoverContent, PopoverRoot, PopoverTitle, PopoverTrigger } from '~/components/ui/popover';
import { Icon, Progress, Spinner, Table, Text } from '@chakra-ui/react';
import { useFileUpload } from '../../hooks/use-file-upload';
import { Tooltip } from '~/components/ui/tooltip';

interface FilePickerProps {
	disabled: boolean;
	identifier: string;
}

export const FilePicker: React.FC<FilePickerProps> = ({
	disabled,
	identifier,
}) => {
	const { data: currentUser } = userApi.useGetOrCreateDefaultUserQuery();
	const { files, uploadFile, removeFileUpload } = useFileUpload(identifier);
	const fileCount = Object.keys(files).length;

	const fileInputRef = useRef<HTMLInputElement>(null);

	const handleFileSelect = () => {
		fileInputRef.current?.click();
	};

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files;
		if (!files || !currentUser) return;

		Array.from(files).forEach(file => {
			uploadFile(file, {
				type: 'user',
				identifier: currentUser!.user.id,
			});
		});

		// Clear the input so the same file can be selected again
		event.target.value = '';
	};

	const renderUploadStatus = (status: string, progress?: number) => {
		switch (status) {
			case 'pending':
				return <Text fontSize={'xs'}>{'Pending...'}</Text>;
			case 'uploading':
				return (
					<Tooltip content={`${progress || 0}%`}>
						<Progress.Root striped animated colorPalette={'blue'} value={progress || 0}>
							<Progress.Track>
								<Progress.Range />
							</Progress.Track>
						</Progress.Root>
					</Tooltip>
				);
			case 'confirming':
				return <Spinner size={'sm'} />;
			case 'ready':
				return <Text fontSize={'xs'} color={'green.500'}>{'Ready'}</Text>;
			case 'error':
				return <Text fontSize={'xs'} color={'red.500'}>{'Error'}</Text>;
			default:
				return null;
		}
	};

	return (
		<PopoverRoot size={'md'}>
			<PopoverTrigger asChild>
				<Button
					disabled={disabled}
					size={'2xs'}
					variant={'outline'}
				>
					<LuFileUp />
					{fileCount > 0 && ` (${fileCount})`}
					<LuChevronDown />
				</Button>
			</PopoverTrigger>
			<PopoverContent>
				<PopoverArrow />
				<PopoverBody py={4}>
					<PopoverTitle fontWeight={'bold'}>
						{'Upload files...'}
					</PopoverTitle>
				</PopoverBody>
				
				<Table.Root width={'calc(100% - 2px)'} ml={'1px'}>
					<Table.Header>
						<Table.Row>
							<Table.ColumnHeader>{'Name'}</Table.ColumnHeader>
							<Table.ColumnHeader>{'Status'}</Table.ColumnHeader>
							<Table.ColumnHeader>{'Actions'}</Table.ColumnHeader>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{Object.keys(files).length > 0 ? (
							Object.entries(files).map(([uploadId, upload]) => (
								<Table.Row key={uploadId}>
									<Table.Cell>
										<Text fontSize="sm" maxW="200px">
											{upload.fileMetadata.name}
										</Text>
									</Table.Cell>
									<Table.Cell>{renderUploadStatus(upload.status, upload.progress)}</Table.Cell>
									<Table.Cell>
										{(upload.status === 'ready' || upload.status === 'error') && (
											<Button 
												size={'2xs'} 
												variant={'surface'} 
												colorPalette={'red'}
												onClick={() => removeFileUpload(uploadId)}
											>
												<LuTrash />
											</Button>
										)}
									</Table.Cell>
								</Table.Row>
							))
						) : (
							<Table.Row>
								<Table.Cell colSpan={3} textAlign="center">
									<Text fontSize={'sm'} color="gray.500">No files uploaded yet</Text>
								</Table.Cell>
							</Table.Row>
						)}
					</Table.Body>
				</Table.Root>

				<PopoverBody>
					<input
						type={'file'}
						ref={fileInputRef}
						onChange={handleFileChange}
						style={{ display: 'none' }}
						multiple
					/>
					<Button
						disabled={disabled}
						size={'2xs'}
						variant={'outline'}
						colorPalette={'blue'}
						w={'full'}
						onClick={handleFileSelect}
					>
						{'Upload file '}
						<Icon>
							<LuFilePlus />
						</Icon>
					</Button>
				</PopoverBody>
			</PopoverContent>
		</PopoverRoot>
	);
}
