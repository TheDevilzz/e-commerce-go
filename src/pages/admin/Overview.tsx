import { useEffect, useState } from "react";
import { DollarSign, ShoppingBag, TrendingUp, Users } from "lucide-react";
import { toast } from "sonner";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { getDashboardStats, type DashboardStats } from "../../api/dashboard";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { formatPrice } from "../../lib/utils";

export function AdminOverviewPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);

  useEffect(() => {
    const loadStats = async () => {
      try {
        const response = await getDashboardStats();
        setStats(response.data);
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load dashboard");
      }
    };
    void loadStats();
  }, []);

  const metricCards = [
    { title: "Total Revenue", value: formatPrice(stats?.total_revenue ?? 0), icon: DollarSign },
    { title: "Orders", value: String(stats?.orders ?? 0), icon: ShoppingBag },
    { title: "Customers", value: String(stats?.customers ?? 0), icon: Users },
    { title: "Avg Order Value", value: formatPrice(stats?.avg_order ?? 0), icon: DollarSign },
  ];

  return (
    <div className="space-y-6">
      <h2>Dashboard Overview</h2>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        {metricCards.map((metric) => {
          const Icon = metric.icon;
          return (
            <Card key={metric.title}>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">{metric.title}</CardTitle>
                <Icon className="size-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metric.value}</div>
                <p className="mt-1 flex items-center gap-1 text-xs text-success">
                  <TrendingUp className="size-3" />
                  Live API data
                </p>
              </CardContent>
            </Card>
          );
        })}
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Sales by Category</CardTitle>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={stats?.category_sales ?? []}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="category" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="sales" fill="hsl(var(--primary))" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
    </div>
  );
}
