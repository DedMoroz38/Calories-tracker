"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, BarChart2, PlusCircle, Images, User } from "lucide-react";
import styles from "./TabBar.module.css";

const LEFT_TABS = [
  { href: "/home",  icon: Home,      label: "Home" },
  { href: "/stats", icon: BarChart2, label: "Stats" },
] as const;

const RIGHT_TABS = [
  { href: "/feed",    icon: Images, label: "Feed" },
  { href: "/profile", icon: User,   label: "Profile" },
] as const;

interface TabBarProps {
  /** Center action button. Adds food on Home/Stats, posts a photo on Feed/Profile. */
  onFabClick: () => void;
  /** Accessible label for the FAB, since its action is page-specific. */
  fabLabel?: string;
}

export function TabBar({ onFabClick, fabLabel = "Add" }: TabBarProps) {
  const pathname = usePathname();

  return (
    <nav className={styles.bar}>
      <div className={styles.inner}>
        {LEFT_TABS.map((t) => (
          <TabItem key={t.href} {...t} active={pathname === t.href} />
        ))}

        <button className={styles.fab} onClick={onFabClick} aria-label={fabLabel}>
          <PlusCircle size={28} strokeWidth={1.8} />
        </button>

        {RIGHT_TABS.map((t) => (
          <TabItem key={t.href} {...t} active={pathname === t.href} />
        ))}
      </div>
    </nav>
  );
}

function TabItem({
  href,
  icon: Icon,
  label,
  active,
}: {
  href: string;
  icon: React.ElementType;
  label: string;
  active: boolean;
}) {
  return (
    <Link
      href={href}
      className={[styles.tab, active ? styles.activeTab : ""].filter(Boolean).join(" ")}
      aria-current={active ? "page" : undefined}
    >
      <Icon size={22} strokeWidth={active ? 2.2 : 1.8} />
      <span className={styles.tabLabel}>{label}</span>
    </Link>
  );
}
