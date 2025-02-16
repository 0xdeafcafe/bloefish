import { motion } from 'motion/react';

export const PanelRoot: React.FC<React.PropsWithChildren> = ({ children }) => (
	<motion.div
		style={{ width: '100%', height: '100%' }}
		initial={{ scale: 0.95, opacity: 0, y: 10 }}
		animate={{ scale: 1, opacity: 1, y: 0 }}
		exit={{ scale: 0.95, opacity: 0, y: 10 }}
		transition={{ duration: 0.3 }}
	>
		{children}
	</motion.div>
);
