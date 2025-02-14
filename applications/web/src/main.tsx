import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router';
import { App } from './App';
import { Welcome } from './welcome/welcome';
import { Provider as ChakraProvider } from './components/ui/provider';

const root = createRoot(document.getElementById('root')!);

root.render(
	<React.StrictMode>
		<ChakraProvider>
			<BrowserRouter>
				<App>
					<Routes>
						<Route path="/" element={<h1>home</h1>} />
					</Routes>
				</App>
			</BrowserRouter>
		</ChakraProvider>
	</React.StrictMode>
);
