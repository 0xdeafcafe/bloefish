import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router';
import { App } from './App';
import { Welcome } from './welcome/welcome';
import { Provider as ChakraProvider } from './components/ui/provider';
import { Provider as ReduxProvider } from 'react-redux';
import { store } from './store';
import { EnsureReadiness } from './components/molecules/EnsureReadiness';

const root = createRoot(document.getElementById('root')!);

root.render(
	<React.StrictMode>
		<ChakraProvider>
			<BrowserRouter>
				<ReduxProvider store={store}>
					<App>
						<EnsureReadiness>
							<Routes>
								<Route path="/" element={<h1>home</h1>} />
								<Route path="testing" element={<Welcome />} />
							</Routes>
						</EnsureReadiness>
					</App>
				</ReduxProvider>
			</BrowserRouter>
		</ChakraProvider>
	</React.StrictMode>
);
