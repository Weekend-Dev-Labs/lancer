/* eslint-disable react/prop-types */

import React, { ReactNode } from "react";

type ContainerProps = {
  children: ReactNode;
  className?: string;
} & React.HTMLAttributes<HTMLDivElement>;

const Container: React.FC<ContainerProps> = ({
  children,
  className = "",
  ...props
}) => {
  return (
    <div
      className={`container w-full h-full xl:max-w-[96rem] mx-auto px-4 lg:px-5 overflow-hidden ${className}`}
      {...props}
    >
      {children}
    </div>
  );
};

export default Container;
