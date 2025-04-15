'use client';
import React, { useState } from 'react';
import Link from 'next/link';
import { ApiToggle } from '@/components/ApiToggle';
import {
  BarChart3,
  CircleDollarSign,
  LineChart,
  Star,
  Cpu,
  Settings,
  MessageSquare,
  Send,
  Zap,
  FileBarChart
} from 'lucide-react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from "@/components/ui/sidebar";

type NavItem = {
  name: string;
  path: string;
  icon: React.ElementType;
};

const navItems: NavItem[] = [
  { name: 'Dashboard', path: '/', icon: BarChart3 },
  { name: 'Portfolio', path: '/portfolio', icon: CircleDollarSign },
  { name: 'Trading', path: '/trading', icon: LineChart },
  { name: 'New Coins', path: '/new-coins', icon: Star },
  { name: 'Backtesting', path: '/backtesting', icon: FileBarChart },
  { name: 'System Status', path: '/system', icon: Cpu },
  { name: 'Bot Config', path: '/config', icon: Settings },
  { name: 'Settings', path: '/settings', icon: Zap },
];

type Message = {
  text: string;
  fromBot: boolean;
  timestamp: Date;
};

export function AppSidebar() {
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState<Message[]>([
    {
      text: "Welcome to BruteBot. How can I assist you today?",
      fromBot: true,
      timestamp: new Date(),
    }
  ]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!message.trim()) return;
    const userMessage: Message = {
      text: message,
      fromBot: false,
      timestamp: new Date(),
    };
    setMessages((prev) => [...prev, userMessage]);
    setMessage('');
    setTimeout(() => {
      const botMessage: Message = {
        text: "I'm analyzing the market. Will update you soon.",
        fromBot: true,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, botMessage]);
    }, 1000);
  };

  return (
    <Sidebar>
      <SidebarHeader className="p-0 border-b border-brutal-border">
        <div className="p-6">
          <h1 className="text-lg font-bold uppercase tracking-widest flex items-center">
            BRUTE<span className="text-sidebar-foreground/30 mx-2">/</span>DASH
          </h1>
        </div>
      </SidebarHeader>
      <SidebarContent className="p-0">
        <nav className="flex-1 py-6">
          <SidebarMenu>
            {navItems.map((item) => (
              <SidebarMenuItem key={item.name}>
                <SidebarMenuButton asChild>
                  <Link href={item.path} className="flex items-center px-6 py-3 text-sm">
                    <item.icon className="h-5 w-5 mr-3" />
                    {item.name}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </nav>
      </SidebarContent>
      <SidebarFooter className="p-3 border-t border-brutal-border">
        <div className="text-sm font-bold flex items-center mb-2 text-brutal-text">
          <MessageSquare className="h-4 w-4 mr-2" />
          BRUTEBOT CHAT
        </div>
        <div className="h-48 overflow-y-auto mb-2 p-2 bg-brutal-panel/30 rounded text-xs">
          {messages.map((msg, index) => (
            <div
              key={index}
              className={`mb-2 ${msg.fromBot ? 'text-brutal-info' : 'text-brutal-text'}`}
            >
              <span className="opacity-70 text-[10px]">
                {msg.fromBot ? 'BOT' : 'YOU'} • {msg.timestamp.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
              </span>
              <div className="mt-1">{msg.text}</div>
            </div>
          ))}
        </div>
        <form onSubmit={handleSendMessage} className="flex items-center">
          <input
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Message BruteBot..."
            className="bg-brutal-panel text-brutal-text text-xs flex-1 rounded-l p-2 outline-none border-l border-y border-brutal-border focus:border-brutal-info"
          />
          <button
            type="submit"
            className="bg-brutal-panel text-brutal-text p-2 rounded-r border-r border-y border-brutal-border hover:bg-brutal-info hover:text-brutal-background"
          >
            <Send className="h-4 w-4" />
          </button>
        </form>
        <div className="mt-3 flex justify-between items-center">
          <div className="text-xs text-brutal-text/50">
            Bot version: 1.4.2
          </div>
          <ApiToggle />
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
