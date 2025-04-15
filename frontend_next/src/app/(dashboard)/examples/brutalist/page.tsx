"use client"

import { BrutalistButton } from "@/components/ui/brutalist-button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export default function BrutalistDemoPage() {
  return (
    <div className="container py-10">
      <h1 className="text-4xl font-extrabold mb-6">Brutalist Components</h1>
      <p className="text-lg mb-8">
        A collection of brutalist-inspired UI components with bold borders, sharp corners, and dynamic interactive states.
      </p>

      <Tabs defaultValue="buttons">
        <TabsList className="mb-8">
          <TabsTrigger value="buttons">Buttons</TabsTrigger>
          <TabsTrigger value="cards">Cards</TabsTrigger>
          <TabsTrigger value="forms">Form Inputs</TabsTrigger>
        </TabsList>

        <TabsContent value="buttons" className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>Button Variants</CardTitle>
              <CardDescription>
                Different button styles for various actions and states.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-4">
              <BrutalistButton>Default</BrutalistButton>
              <BrutalistButton variant="outline">Outline</BrutalistButton>
              <BrutalistButton variant="danger">Danger</BrutalistButton>
              <BrutalistButton variant="success">Success</BrutalistButton>
              <BrutalistButton variant="warning">Warning</BrutalistButton>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Button Sizes</CardTitle>
              <CardDescription>
                Different size options for buttons to fit various layouts.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap items-center gap-4">
              <BrutalistButton size="sm">Small</BrutalistButton>
              <BrutalistButton>Default</BrutalistButton>
              <BrutalistButton size="lg">Large</BrutalistButton>
              <BrutalistButton size="xl">Extra Large</BrutalistButton>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Button Shadows</CardTitle>
              <CardDescription>
                Shadow variations for depth and emphasis.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-4">
              <BrutalistButton shadow="none">No Shadow</BrutalistButton>
              <BrutalistButton shadow="sm">Small Shadow</BrutalistButton>
              <BrutalistButton>Default Shadow</BrutalistButton>
              <BrutalistButton shadow="lg">Large Shadow</BrutalistButton>
              <BrutalistButton shadow="xl">Extra Large Shadow</BrutalistButton>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Transform Effects</CardTitle>
              <CardDescription>
                Interactive transform effects on hover and click.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-4">
              <BrutalistButton transform="none">No Transform</BrutalistButton>
              <BrutalistButton transform="sm">Small Transform</BrutalistButton>
              <BrutalistButton>Default Transform</BrutalistButton>
              <BrutalistButton transform="lg">Large Transform</BrutalistButton>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Border Widths</CardTitle>
              <CardDescription>
                Different border thickness options.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-4">
              <BrutalistButton borderWidth="none">No Border</BrutalistButton>
              <BrutalistButton borderWidth="sm">Thin Border</BrutalistButton>
              <BrutalistButton>Default Border</BrutalistButton>
              <BrutalistButton borderWidth="lg">Thick Border</BrutalistButton>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Combined Styles</CardTitle>
              <CardDescription>
                Combining different props for customized buttons.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-4">
              <BrutalistButton
                variant="danger"
                size="lg"
                shadow="xl"
                borderWidth="lg"
              >
                Delete Account
              </BrutalistButton>
              
              <BrutalistButton
                variant="success"
                size="lg"
                transform="lg"
              >
                Confirm Payment
              </BrutalistButton>
              
              <BrutalistButton
                variant="outline"
                size="sm"
                shadow="sm"
                borderWidth="sm"
              >
                Cancel
              </BrutalistButton>
              
              <BrutalistButton
                variant="warning"
                size="xl"
                shadow="lg"
                transform="sm"
              >
                Review Changes
              </BrutalistButton>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="cards" className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="brutalist-card">
            <h3 className="text-xl font-bold mb-2">Default Card</h3>
            <p>A basic brutalist card with default styling.</p>
            <div className="mt-4">
              <BrutalistButton size="sm">Action</BrutalistButton>
            </div>
          </div>

          <div className="brutalist-card danger">
            <h3 className="text-xl font-bold mb-2">Danger Card</h3>
            <p>A card with danger styling for warnings and critical information.</p>
            <div className="mt-4">
              <BrutalistButton variant="danger" size="sm">Delete</BrutalistButton>
            </div>
          </div>

          <div className="brutalist-card success">
            <h3 className="text-xl font-bold mb-2">Success Card</h3>
            <p>A card with success styling for positive outcomes and confirmations.</p>
            <div className="mt-4">
              <BrutalistButton variant="success" size="sm">Confirm</BrutalistButton>
            </div>
          </div>

          <div className="brutalist-card warning">
            <h3 className="text-xl font-bold mb-2">Warning Card</h3>
            <p>A card with warning styling for important notices that require attention.</p>
            <div className="mt-4">
              <BrutalistButton variant="warning" size="sm">Review</BrutalistButton>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="forms" className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>Brutalist Form Elements</CardTitle>
              <CardDescription>
                Form controls with brutalist styling.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="block font-bold mb-2">Text Input</label>
                <input type="text" placeholder="Enter your name" className="brutalist-input" />
              </div>
              
              <div>
                <label className="block font-bold mb-2">Select Input</label>
                <select className="brutalist-select w-full">
                  <option>Bitcoin</option>
                  <option>Ethereum</option>
                  <option>Solana</option>
                  <option>Cardano</option>
                </select>
              </div>
              
              <div>
                <label className="block font-bold mb-2">Textarea</label>
                <textarea className="brutalist-input" rows={4} placeholder="Enter your message"></textarea>
              </div>
            </CardContent>
            <CardFooter>
              <BrutalistButton>Submit Form</BrutalistButton>
            </CardFooter>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
} 