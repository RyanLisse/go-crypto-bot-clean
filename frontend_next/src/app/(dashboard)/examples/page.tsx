"use client"

import React from "react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import Link from "next/link"
import { ArrowRight } from "lucide-react"

export default function ExamplesPage() {
  const examples = [
    {
      title: "Brutalist Button",
      description: "A demo of the Brutalist Button component with all variants and sizes",
      href: "/examples/brutalist-button",
    },
    // Add more examples here as they are created
  ]

  return (
    <div className="container py-10">
      <h1 className="text-4xl font-bold mb-2">Component Examples</h1>
      <p className="text-muted-foreground mb-8">
        Browse through component examples and demos to learn how to use them in your project.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {examples.map((example, index) => (
          <Link href={example.href} key={index} className="block">
            <Card className="h-full transition-all hover:shadow-md">
              <CardHeader>
                <CardTitle>{example.title}</CardTitle>
                <CardDescription>{example.description}</CardDescription>
              </CardHeader>
              <CardFooter className="flex justify-end">
                <div className="flex items-center text-sm">
                  View example <ArrowRight className="ml-2 h-4 w-4" />
                </div>
              </CardFooter>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
} 