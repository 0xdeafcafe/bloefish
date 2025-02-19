import { IconButton, Status } from "@chakra-ui/react";
import { motion } from "motion/react";
import React from "react";
import { useEffect, useState } from "react";
import { Tooltip } from "../ui/tooltip";

type Status = 'online' | 'offline' | 'pending';

const serviceHealthPaths: string[] = [
	'http://svc_ai_relay.bloefish.local:4003/system/health',
	'http://svc_conversation.bloefish.local:4002/system/health',
	'http://svc_user.bloefish.local:4001/system/health',
	'http://svc_stream.bloefish.local:4004/system/health',
	'http://svc_file_upload.bloefish.local:4005/system/health',
];

export const PlatformStatus: React.FC = () => {
	const [status, setStatus] = useState<Status>('pending');
	const [refreshing, setRefreshing] = useState(false);

	async function checkStatus() {
		if (refreshing) return;
		
		try {
			setRefreshing(true);
			await Promise.all(serviceHealthPaths.map(async (svcPath) => fetch(svcPath)));
			setStatus('online');
		} catch {
			setStatus('offline');
		} finally {
			setRefreshing(false);
		}
	}

	useEffect(() => {
		checkStatus();

		const interval = window.setInterval(checkStatus, 60000);

		return () => {
			window.clearInterval(interval);
		};
	}, []);

	switch (true) {
		case refreshing:
			return (
				<Status.Root colorPalette={'orange'} size={'sm'}>
					<MotionIndicator
						animate={{ opacity: [0.5, 1, 0.5] }}
						transition={{ duration: 1.5, repeat: Infinity, ease: "easeInOut" }}
					/>
					{'Checking platform status...'}
				</Status.Root>
			);

		case status === 'pending':
			return (
				<Tooltip content="...">
					<Status.Root colorPalette={'orange'} size={'sm'}>
						<MotionIndicator
							animate={{ opacity: [0.5, 1, 0.5] }}
							transition={{ duration: 1.5, repeat: Infinity, ease: "easeInOut" }}
						/>
						{'Platform status pending'}
					</Status.Root>
				</Tooltip>
			);

		case status === 'offline':
			return (
				<Tooltip content={'Bloefish platform is offline. You can click to try again.'}>
					<Status.Root colorPalette={'red'} size={'sm'}>
						<Status.Indicator />
						{'Offline'}
					</Status.Root>
				</Tooltip>
			);

		case status === 'online':
			return (
				<Tooltip content={'Bloefish platform is online'}>
					<Status.Root colorPalette={'green'} size={'sm'}>
						<Status.Indicator />
						{'Online'}
					</Status.Root>
				</Tooltip>
			);

		default: return null;
	}
}

const MotionIndicator = motion(Status.Indicator);
