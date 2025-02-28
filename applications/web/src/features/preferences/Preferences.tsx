import { Container, Fieldset, Grid, GridItem, Input, Stack, Switch, Textarea } from '@chakra-ui/react';
import React from 'react';
import { BorderedScrollContainer } from '~/components/atoms/BorderedScrollContainer';
import { HeaderCard } from '~/components/atoms/HeaderCard';
import { Button } from '~/components/ui/button';
import { Field } from '~/components/ui/field';
import { NativeSelectField, NativeSelectRoot } from '~/components/ui/native-select';

export const Preferences: React.FC = () => (
	<Grid
		templateRows={'auto auto 1fr'}
		maxH={'100%'}
		overflow={'hidden'}
	>
		<GridItem p={6}>
			<HeaderCard
				title={'Preferences'}
				description={'Make it yours, or if you have terrible taste maybe it\'s best to leave it as it is honestly.'}
			/>
		</GridItem>

		<GridItem asChild>
			<BorderedScrollContainer p={6} triggerOffset={24}>
				<Container maxW={'4xl'}>
					<Stack gap={16}>
						<Fieldset.Root
							size="lg"
							display={'flex'}
							flexDirection={'row'}
							gap={10}
						>
							<Stack flex={1}>
								<Fieldset.Legend fontSize={'xl'}>
									{'Profile details'}
								</Fieldset.Legend>
								<Fieldset.HelperText>
									{'Who are you? You can lie, nobody will notice.'}
								</Fieldset.HelperText>
							</Stack>

							<Fieldset.Content flex={1} mt={0}>
								<Field label="Name">
									<Input name="name" disabled value={'Alexander Forbes-Reed'} />
								</Field>

								<Field label="Bio">
									<Textarea name="bio" disabled autoresize placeholder={'Write a bit about yourself, your interests, etc.'} />
								</Field>

								<Button disabled type={'submit'} variant={'surface'} alignSelf={'flex-start'}>
									{'Update'}
								</Button>
							</Fieldset.Content>
						</Fieldset.Root>

						<Fieldset.Root
							size="lg"
							display={'flex'}
							flexDirection={'row'}
							gap={10}
						>
							<Stack flex={1}>
								<Fieldset.Legend fontSize={'xl'}>
									{'Assistant preferences'}
								</Fieldset.Legend>
								<Fieldset.HelperText>
									{'Tell your AI assistant how to behave (it\'ll probably ignore you anyway)'}
								</Fieldset.HelperText>
							</Stack>

							<Fieldset.Content flex={1} mt={0}>
								<Field label="Personal prompt">
									<Textarea
										name="personalPrompt"
										autoresize
										disabled
										placeholder={'Write some custom instructions to guide your assistant for every message you send'}
									/>
								</Field>

								<Field label="Assistant language">
									<NativeSelectRoot>
										<NativeSelectField
											name={'assistantLanguage'}
											disabled
											defaultValue={'English (UK)'}
											items={[
												'English (UK)',
												'English (US)',
												'English (AU)',
												'English (CA)',
												'English (NZ)',
												'Nederlands',
												'Français',
												'Deutsch',
												'Español',
												'Italiano',
												'Português',
												'Polski',
												'Русский',
											]}
										/>
									</NativeSelectRoot>
								</Field>

								<Field label="Title generation" gap={2}>
									<Switch.Root checked={false} disabled>
										<Switch.HiddenInput />
										<Switch.Control>
											<Switch.Thumb />
										</Switch.Control>
										<Switch.Label>
											{'Use separate model for title generation than for messages'}
										</Switch.Label>
									</Switch.Root>

									<NativeSelectRoot>
										<NativeSelectField
											name={'titleGenerationModel'}
											disabled
											defaultValue={'Open AI (GPT 4o mini)'}
											items={[
												'Open AI (GPT 4o mini)',
											]}
										/>
									</NativeSelectRoot>
								</Field>

								<Button disabled type={'submit'} variant={'surface'} alignSelf={'flex-start'}>
									{'Update'}
								</Button>
							</Fieldset.Content>
						</Fieldset.Root>

						<Fieldset.Root
							size="lg"
							display={'flex'}
							flexDirection={'row'}
							gap={10}
						>
							<Stack flex={1}>
								<Fieldset.Legend fontSize={'xl'}>
									{'Display'}
								</Fieldset.Legend>
								<Fieldset.HelperText>
									{"Make your eyes hate you less (or more, we don't judge)"}
								</Fieldset.HelperText>
							</Stack>

							<Fieldset.Content flex={1} mt={0}>
								<Field label={'Theme'}>
									<NativeSelectRoot>
										<NativeSelectField
											name={'theme'}
											disabled
											defaultValue={'System'}
											items={[
												'Light',
												'Dark',
												'System',
											]}
										/>
									</NativeSelectRoot>
								</Field>

								<Field label={'Reduced motion'}>
									<Switch.Root disabled>
										<Switch.HiddenInput />
										<Switch.Control>
											<Switch.Thumb />
										</Switch.Control>
										<Switch.Label>
											{'Normal'}
										</Switch.Label>
									</Switch.Root>
								</Field>

								<Field label={'Background animations'}>
									<Switch.Root checked disabled>
										<Switch.HiddenInput />
										<Switch.Control>
											<Switch.Thumb />
										</Switch.Control>
										<Switch.Label>
											{'Enabled'}
										</Switch.Label>
									</Switch.Root>
								</Field>
							</Fieldset.Content>
						</Fieldset.Root>
					</Stack>
				</Container>
			</BorderedScrollContainer>
		</GridItem>
	</Grid>
);
