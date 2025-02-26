import { Button, Fieldset, Input, Stack, Textarea } from '@chakra-ui/react';
import { DialogActionTrigger, DialogBody, DialogCloseTrigger, DialogContent, DialogFooter, DialogHeader, DialogRoot, DialogTitle, DialogTrigger } from '~/components/ui/dialog';
import { useRef } from 'react';
import { skillSetApi } from '~/api/bloefish/skill-set';
import { Field } from '~/components/ui/field';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { userApi } from '~/api/bloefish/user';
import { toaster } from '~/components/ui/toaster';

const schema = z.object({
	name: z.string().min(2, 'Name must be at least 2 characters'),
	description: z.string().min(2, 'Description must be at least 2 characters'),
	prompt: z.string()
		.min(2, 'Instructions must be at least 2 characters')
		.max(3000, 'Instructions must be at most 3000 characters'),
});

type FormData = z.infer<typeof schema>;

interface CreateSkillSetDialogProps {
	onSuccess?: () => void;
}

export const CreateSkillSetDialog: React.FC<React.PropsWithChildren<CreateSkillSetDialogProps>> = ({
	children,
	onSuccess,
}) => {
	const { data: currentUser } = userApi.useGetOrCreateDefaultUserQuery();
	const initialFocusRef = useRef<HTMLInputElement>(null)
	const closeRef = useRef<HTMLButtonElement>(null);
	const [createSkillSet] = skillSetApi.useCreateSkillSetMutation();

	const {
		register,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm({
		resolver: zodResolver(schema),
	});

	const onSubmit = async (data: FormData) => {
		try {
			await createSkillSet({
				name: data.name,
				description: data.description,
				prompt: data.prompt,
				icon: 'mic_vocal',
				owner: {
					type: 'user',
					identifier: currentUser!.user.id,
				},
			}).unwrap();

			closeRef.current?.click();

			onSuccess?.();

			toaster.create({
				title: 'Skill set created',
				description: 'The skill set has been created successfully',
				type: 'success',
			});
		} catch (error) {
			console.error(error);
			toaster.create({
				title: 'Failed to create skill set',
				description: 'An error occurred while creating the skill set',
				type: 'error',
			})
		}
	};

	return (
		<DialogRoot size={'lg'} initialFocusEl={() => initialFocusRef.current}>
			<DialogTrigger asChild>
				{children}
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Create new skill set</DialogTitle>
				</DialogHeader>
				<DialogBody>
					{/* eslint-disable-next-line react-compiler/react-compiler */}
					<form id={'create-skill-set-form'} onSubmit={handleSubmit(onSubmit)}>
						<Fieldset.Root>
							<Stack>
								<Fieldset.Legend>Back to school</Fieldset.Legend>
								<Fieldset.HelperText>
									{'Skill sets store context about subjects, skills, '}
									{'or instructions that can be used to guide AI '}
									{'conversations.'}
								</Fieldset.HelperText>
							</Stack>

							<Fieldset.Content>
								<Field label="Name" required invalid={Boolean(errors.name)}>
									<Input {...register('name')} type={'text'} />
								</Field>

								<Field label="Description" required invalid={Boolean(errors.description)}>
									<Input {...register('description')} type={'text'} />
								</Field>

								{/* <Field
									label={'Icon'}
									invalid={Boolean(errors.icon)}
								>
									<Controller
										control={control}
										name={'icon'}
										render={({ field }) => (
											<SelectRoot
												name={field.name}
												value={field.value}
												onValueChange={({ value }) => field.onChange(value)}
												onInteractOutside={() => field.onBlur()}
												collection={getLucideIconNames()}
											>
												<SelectTrigger>
													<SelectValueText placeholder={'Select icon'} />
												</SelectTrigger>
												<SelectContent>
													{getLucideIconNames().map((icon) => (
														<SelectItem item={movie} key={movie.value}>
															{movie.label}
														</SelectItem>
													))}
												</SelectContent>
											</SelectRoot>
										)}
									/>
								</Field> */}

								<Field label="Instructions" required invalid={Boolean(errors.prompt)}>
									<Textarea
										{...register('prompt')}
										name="prompt"
										autoresize
										maxH={'40'}
										minH={'5'}
										placeholder={'Instructions that you want your assistant to follow'}
										variant={'outline'}
									/>
								</Field>
							</Fieldset.Content>
						</Fieldset.Root>
					</form>
				</DialogBody>
				<DialogFooter>
					<DialogActionTrigger asChild>
						<Button
							disabled={isSubmitting}
							variant={'outline'}
							size={'sm'}
						>
							Cancel
						</Button>
					</DialogActionTrigger>
					<Button
						type={'submit'}
						variant={'surface'}
						colorPalette={'blue'}
						size={'sm'}
						loading={isSubmitting}
						form="create-skill-set-form"
					>
						Create
					</Button>
				</DialogFooter>
				<DialogCloseTrigger ref={closeRef} />
			</DialogContent>
		</DialogRoot>
	);
};
