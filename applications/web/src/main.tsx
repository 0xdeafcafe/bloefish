import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router';
import { App } from './App';
import { Provider as ChakraProvider } from './components/ui/provider';
import { Provider as ReduxProvider } from 'react-redux';
import { store } from './store';
import { EnsureReadiness } from './components/molecules/EnsureReadiness';
import { NotFound } from './pages/NotFound';
import { PanelRoot } from './components/molecules/PanelRoot';
import { NewConversation } from './features/new-conversation/NewConversation';
import { Theme } from '@chakra-ui/react';
import { Conversation } from './features/conversations/Conversation';
import { Ready } from './components/molecules/Ready';
import { HelmetProvider } from 'react-helmet-async';
import { ConversationsList } from './features/conversations/ConversationsList';
import { OmniBar } from './features/omnibar/OmniBar';
import { Toaster } from './components/ui/toaster';

const root = createRoot(document.getElementById('root')!);


root.render(
	<React.StrictMode>
		<HelmetProvider>
			<ChakraProvider>
				<BrowserRouter>
					<Theme>
						<ReduxProvider store={store}>
							<App>
								<EnsureReadiness>
									<Ready>
										<Routes>
											<Route path="/" element={wrap(NewConversation)} />
											<Route path="/conversations/:conversationId" element={wrap(Conversation)} />
											<Route path="/conversations" element={wrap(ConversationsList)} />

											<Route path="*" element={wrap(NotFound)} />
										</Routes>
									</Ready>
								</EnsureReadiness>
							</App>
							<Toaster />
							<OmniBar />
						</ReduxProvider>
					</Theme>
				</BrowserRouter>
			</ChakraProvider>
		</HelmetProvider>
	</React.StrictMode>
);

function wrap(Component: React.FC) {
	return (
		<PanelRoot>
			{<Component />}
		</PanelRoot>
	)
}
