import { conversationApi } from '~/api/bloefish/conversation';
import { userApi } from '~/api/bloefish/user';

export const Ready: React.FC<React.PropsWithChildren> = ({ children }) => {
	const { data: userData } = userApi.useGetOrCreateDefaultUserQuery();
	
	conversationApi.useListConversationsWithInteractionsQuery({
		owner: {
			type: 'user',
			identifier: userData?.user.id ?? 'impossible',
		},
	}, {
		skip: !userData, // should never happen!
	});

	return children;
};
