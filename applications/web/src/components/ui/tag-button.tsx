import { Tag as ChakraTag } from "@chakra-ui/react"
import * as React from "react"

export interface TagButtonProps extends ChakraTag.RootProps {
  onClicked: (prompt: string) => void;
}

export const TagButton = React.forwardRef<HTMLSpanElement, TagButtonProps>(
  function TagButton(props, ref) {
    const {
      children,
      onClicked,
      ...rest
    } = props

    return (
      <ChakraTag.Root
        tabIndex={0}
        ref={ref}
        px={4}
        py={2}
        cursor={'pointer'}
        onClick={() => onClicked(children as string)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === 'Space') onClicked(children as string);
        }}
        {...rest}
      >
        <ChakraTag.Label>{children}</ChakraTag.Label>
      </ChakraTag.Root>
    )
  },
)
