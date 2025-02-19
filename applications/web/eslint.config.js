import globals from 'globals';
import pluginJs from '@eslint/js';
import tseslint from 'typescript-eslint';
import pluginReact from 'eslint-plugin-react';
import reactCompiler from 'eslint-plugin-react-compiler';

/** @type {import('eslint').Linter.Config[]} */
export default [
	{
		ignores: [
			'**/dist/*',
			'tsconfig.json',
		]
	},
	{ files: ['src/**/*.{ts,tsx}'] },
	{ languageOptions: { globals: globals.browser } },
	pluginJs.configs.recommended,
	...tseslint.configs.recommended,
	pluginReact.configs.flat.recommended,
	{
		plugins: { 'react-compiler': reactCompiler },
		rules: { 'react-compiler/react-compiler': 'error' },
	},
	{
		rules: {
			'no-extra-boolean-cast': 'off',

			'react/react-in-jsx-scope': 'off',
			'react/prop-types': 'off',
		},
	}
];
