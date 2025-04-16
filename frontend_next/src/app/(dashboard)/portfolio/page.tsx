

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { cn } from "@/lib/utils";

import React, { useEffect, useState } from "react";

export default function Portfolio() {
  const [wallet, setWallet] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Example upcoming coins data (still mock)
  const upcoming = [
    { symbol: "META", name: "MetaChain", release: "Apr 15", potential: 5, forecast: "+12.7%", forecastClass: "text-green-500" },
    { symbol: "QNT", name: "QuantumNet", release: "Apr 18", potential: 5, forecast: "+4.2%", forecastClass: "text-green-500" },
    { symbol: "AIX", name: "AI Exchange", release: "Apr 21", potential: 5, forecast: "+1.9%", forecastClass: "text-green-500" },
    { symbol: "DFX", name: "DeFi X", release: "Apr 24", potential: 3, forecast: "-2.1%", forecastClass: "text-red-500" },
    { symbol: "WEB4", name: "Web4Token", release: "Apr 29", potential: 4, forecast: "+8.3%", forecastClass: "text-green-500" },
  ];

  useEffect(() => {
    setLoading(true);
    fetch("/api/v1/account/wallet")
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch account data");
        return res.json();
      })
      .then((data) => {
        setWallet(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message || "Unknown error");
        setLoading(false);
      });
  }, []);

  // Prepare summary data from wallet
  const summary = [
    {
      label: "PORTFOLIO VALUE",
      value: wallet ? `$${Number(wallet.totalUSDValue || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : "-",
      sub: undefined,
      subClass: "",
    },
    { label: "ACTIVE TRADES", value: "0" },
    { label: "WIN RATE", value: "0.0%" },
    { label: "AVG PROFIT/TRADE", value: "$0.00" },
  ];

  return (
    <div className="w-full px-4 py-8 max-w-[1400px] mx-auto font-mono">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        {summary.map((item, i) => (
          <Card key={item.label} className="bg-background border-2 border-foreground rounded-none shadow-none">
            <CardContent className="p-4 flex flex-col gap-1">
              <div className="text-xs text-muted-foreground tracking-widest">{item.label}</div>
              <div className="text-xl font-bold">{item.value}</div>
              {item.sub && <div className={cn("text-xs font-semibold", item.subClass)}>{item.sub}</div>}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Portfolio Performance + Account Balance */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        {/* Performance Chart Placeholder */}
        <Card className="col-span-2 bg-background border-2 border-foreground rounded-none">
          <CardHeader>
            <CardTitle className="text-base">PORTFOLIO PERFORMANCE</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-56 flex items-center justify-center border border-dashed border-muted-foreground rounded-none bg-muted">
              <span className="text-muted-foreground">[Chart Placeholder]</span>
            </div>
          </CardContent>
        </Card>
        {/* Account Balance */}
        <Card className="bg-background border-2 border-foreground rounded-none">
          <CardHeader>
            <CardTitle className="text-base">ACCOUNT BALANCE</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="mb-2 text-xs text-muted-foreground flex justify-between">
              <span>Asset</span>
              <span>Amount</span>
              <span>Value (USD)</span>
            </div>
            {loading ? (
              <div className="text-center text-muted-foreground py-4">Loading...</div>
            ) : error ? (
              <div className="text-center text-red-500 py-4">{error}</div>
            ) : wallet && wallet.balances && wallet.balances.length > 0 ? (
              wallet.balances.map((bal: any) => (
                <div key={bal.asset} className="flex justify-between text-sm py-1 border-b border-muted-foreground/20 last:border-b-0">
                  <span>{bal.asset}</span>
                  <span>{bal.free ?? bal.amount ?? "-"}</span>
                  <span>
                    {bal.usdValue !== undefined
                      ? `$${Number(bal.usdValue).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                      : bal.value ?? "-"}
                  </span>
                </div>
              ))
            ) : (
              <div className="text-center text-muted-foreground py-4">No balances found.</div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Upcoming Coins Table */}
      <Card className="bg-background border-2 border-foreground rounded-none">
        <CardHeader>
          <CardTitle className="text-base tracking-wide">UPCOMING COINS (TODAY & TOMORROW)</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow className="bg-muted text-muted-foreground text-xs uppercase">
                <TableHead className="w-1/12">Symbol</TableHead>
                <TableHead className="w-3/12">Name</TableHead>
                <TableHead className="w-2/12">Release</TableHead>
                <TableHead className="w-2/12">Potential</TableHead>
                <TableHead className="w-2/12">Forecast</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {upcoming.map((coin) => (
                <TableRow key={coin.symbol}>
                  <TableCell className="font-bold text-blue-400">{coin.symbol}</TableCell>
                  <TableCell>{coin.name}</TableCell>
                  <TableCell>{coin.release}</TableCell>
                  <TableCell>{"★".repeat(coin.potential)}{"☆".repeat(5 - coin.potential)}</TableCell>
                  <TableCell className={coin.forecastClass}>{coin.forecast}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
