import { LuChevronDown, LuGraduationCap } from 'react-icons/lu';
import { Button } from '~/components/ui/button';
import { MenuCheckboxItem, MenuContent, MenuItemGroup, MenuRoot, MenuTrigger } from '~/components/ui/menu';
import { useAppDispatch } from '~/store';
import { updateSkillSetIds } from '../../store';
import { useChatInput } from '../../hooks/use-chat-input';
import { userApi } from '~/api/bloefish/user';
import { skillSetApi } from '~/api/bloefish/skill-set';

interface SkillSetPickerProps {
	disabled: boolean;
	identifier: string;
	inputRef: React.RefObject<HTMLTextAreaElement | null>;
}

export const SkillSetPicker: React.FC<SkillSetPickerProps> = ({
	disabled,
	identifier,
	inputRef,
}) => {
	const { data: currentUser } = userApi.useGetOrCreateDefaultUserQuery();
	const {
		data: availableSkillSets,
		isSuccess: hasAvailableSkillSets,
		isFetching: isFetchingSkillSets,
	} = skillSetApi.useListSkillSetsByOwnerQuery({
		owner: {
			type: 'user',
			identifier: currentUser!.user.id,
		},
	});

	const { skillSetIds } = useChatInput(identifier);
	const dispatch = useAppDispatch();

	return (
		<MenuRoot onExitComplete={() => inputRef.current?.focus()} >
			<MenuTrigger asChild>
				<Button
					disabled={disabled || !hasAvailableSkillSets}
					loading={isFetchingSkillSets}
					size={'2xs'}
					variant={'outline'}
				>
					<LuGraduationCap />
					{skillSetIds.length > 0 && `(${skillSetIds.length})`}
					<LuChevronDown />
				</Button>
			</MenuTrigger>
			<MenuContent>
				<MenuItemGroup title={'Skill sets'}>
					{availableSkillSets?.skillSets?.map(ss => (
						<MenuCheckboxItem
							key={ss.id}
							value={ss.id}
							checked={skillSetIds.includes(ss.id)}
							onCheckedChange={(checked) => {
								if (checked) {
									dispatch(updateSkillSetIds({
										identifier,
										skillSetIds: [...skillSetIds, ss.id],
									}));
								} else {
									dispatch(updateSkillSetIds({
										identifier,
										skillSetIds: skillSetIds.filter(id => id !== ss.id),
									}));
								}
							}}
						>
							{ss.name}
						</MenuCheckboxItem>
					))}
				</MenuItemGroup>
			</MenuContent>
		</MenuRoot>
	);
}
